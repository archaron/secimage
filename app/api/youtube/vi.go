package youtube

import (
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/labstack/echo"
	"github.com/nfnt/resize"
	"go.uber.org/zap"
)

type (
	ViRequest struct {
		ID            string `params:"id" validate:"required"`
		Size          string `params:"size" validate:"required"`
		File          string `params:"file" validate:"required"`
		width, height uint64
	}

	ViParams struct {
		log          *zap.SugaredLogger
		cli          *http.Client
		savePath     string
		cachePath    string
		allowedSizes []string
		quality      int
	}
)

var re = regexp.MustCompile(`(\d+)x(\d+)`)

func Vi(p ViParams) (echo.HandlerFunc, error) {
	var allowed = make(map[string]struct{}, len(p.allowedSizes))

	// fill up size hashmap
	for _, size := range p.allowedSizes {
		allowed[size] = struct{}{}
	}

	// creates all needed path, if they're not exists
	for _, item := range []string{p.cachePath, p.savePath} {
		if _, err := os.Stat(item); os.IsNotExist(err) {
			if err := os.MkdirAll(item, 0777); err != nil {
				return nil, err
			}
		}
	}

	return func(ctx echo.Context) error {
		var req ViRequest
		var originalFile string
		var err error
		var img image.Image

		if err := ctx.Bind(&req); err != nil {
			p.log.Errorw("can't bind/validate request",
				"error", err)

			return ctx.String(http.StatusBadRequest, "FAIL")
		}
		match := re.FindAllStringSubmatch(req.Size, -1)

		if len(match) < 0 || len(match[0]) < 3 {
			p.log.Errorw("size not match regexp")
			return ctx.String(http.StatusBadRequest, "FAIL SIZE")
		}

		var (
			size   = match[0][0]
			width  = match[0][1]
			height = match[0][2]
		)

		if _, ok := allowed[size]; !ok {
			p.log.Errorw("size not allowed")
			return ctx.String(http.StatusBadRequest, "Size is not allowed")
		}

		cacheSavePath := p.cachePath + "/" + size
		sizedFile := cacheSavePath + "/" + req.ID + ".jpg"
		originalFile = p.savePath + "/" + req.ID + ".jpg"

		if _, err = os.Stat(sizedFile); os.IsNotExist(err) {
			req.width, err = strconv.ParseUint(width, 10, 64)
			if err != nil {
				return err
			}

			req.height, err = strconv.ParseUint(height, 10, 64)
			if err != nil {
				return err
			}

			// ensure directory exists
			if _, err := os.Stat(cacheSavePath); os.IsNotExist(err) {
				if err := os.MkdirAll(cacheSavePath, 0777); err != nil {
					p.log.Errorw("can't create folder",
						"error", err)
					return ctx.String(http.StatusBadRequest, "FAIL")
				}
			}

			if _, err := os.Stat(originalFile); os.IsNotExist(err) {
				p.log.Debugf("%s | %s | DOWNLOAD", req.ID, size)
				p.log.Debugf("Downloading %s ...", req.ID)

				// Get the data
				hqdefault := "https://i.ytimg.com/vi/" + req.ID + "/hqdefault.jpg"
				resp, err := http.Get(hqdefault)
				if err != nil {
					p.log.Debugf("download error %s, file: %s", err, hqdefault)
					return err
				}
				defer resp.Body.Close()

				out, err := os.Create(originalFile)
				if err != nil {
					p.log.Debugf("CreateFile error %s", err)
					return err
				}
				defer out.Close()

				// Write the body to file
				_, err = io.Copy(out, resp.Body)
				if err != nil {
					p.log.Debugw("coping file fail",
						"error", err)
					return err
				}

			} else {
				p.log.Debugf("%s | %s | CACHED ORIGINAL", req.ID, size)
			}

			in, err := os.Open(originalFile)
			if err != nil {
				return err
			}
			defer in.Close()

			img, _, err = image.Decode(in)

			// Image decode error catch
			if err != nil {
				p.log.Debugf("imageDecode error %s", err)
				return err
			}

			resizedImage := resize.Resize(uint(req.width), uint(req.height), img, resize.MitchellNetravali)

			rout, err := os.Create(sizedFile)
			if err != nil {
				p.log.Debugf("outfileCreate error %s", err)
				return err
			}
			defer rout.Close()

			err = jpeg.Encode(rout, resizedImage, &jpeg.Options{
				Quality: p.quality,
			})

			if err != nil {
				p.log.Debugf("imageEncode error %s", err)
				return err
			}
		}

		return ctx.File(sizedFile)
	}, nil
}

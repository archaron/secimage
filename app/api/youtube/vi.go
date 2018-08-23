package youtube

import (
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"image"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/nfnt/resize"
	"image/jpeg"
)

type ViRequest struct {
}

func Vi(log *zap.SugaredLogger, cli *http.Client, savePath string, cachePath string, allowedSizes []string, quality int) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var req ViRequest
		var originalFile string
		var err error
		var img image.Image

		if err := ctx.Bind(&req); err != nil {
			log.Info(zap.Error(err))
			return ctx.String(http.StatusBadRequest, "FAIL")
		}

		file := ctx.Param("file")
		size := ctx.Param("size")
		id := ctx.Param("id")
		var width, height int64

		var re = regexp.MustCompile(`(\d+)x(\d+)`)

		match := re.FindAllStringSubmatch(size, -1)

		if len(match) < 0 || len(match[0]) < 3 {
			return ctx.String(http.StatusBadRequest, "FAIL SIZE")
		}

		// Check, if size is in allowed sizes list
		var sizeAllowed bool
		for _, a := range allowedSizes {
			if a == match[0][0] {
				sizeAllowed = true
				break
			}
		}

		if !sizeAllowed {
			return ctx.String(http.StatusBadRequest, "Size is not allowed")
		}

		width, err = strconv.ParseInt(match[0][1], 10, 64)
		if err != nil {
			return err
		}

		height, err = strconv.ParseInt(match[0][2], 10, 64)
		if err != nil {
			return err
		}

		cacheSavePath := cachePath + "/" + match[0][0]

		originalFile = savePath + "/" + id + ".jpg"

		// ensure directory exists
		if _, err := os.Stat(cacheSavePath); os.IsNotExist(err) {
			os.MkdirAll(cacheSavePath, 0777)
		}

		if _, err := os.Stat(originalFile); os.IsNotExist(err) {
			log.Debugf("%s | %s | DOWNLOAD", id, size)
			log.Debugf("Downloading %s ...", id)

			// Get the data
			resp, err := http.Get("https://i.ytimg.com/vi/" + id + "/hqdefault.jpg")
			if err != nil {
				log.Debugf("download error %s, file: %s", err, "https://i.ytimg.com/vi/"+id+"/hqdefault.jpg")
				return err
			}
			defer resp.Body.Close()

			out, err := os.Create(originalFile)
			if err != nil {
				log.Debugf("CreateFile error %s", err)
				return err
			}
			defer out.Close()

			// Write the body to file
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				log.Debugf("CopyFile error %s", err)
				return err
			}

		} else {
			log.Debugf("%s | %s | CACHED ORIGINAL", id, size)
		}

		in, err := os.Open(originalFile)
		if err != nil {
			return err
		}
		defer in.Close()

		img, _, err = image.Decode(in)

		// Image decode error catch
		if err != nil {
			log.Debugf("imageDecode error %s", err)
			return err
		}

		resizedImage := resize.Resize(uint(width), uint(height), img, resize.MitchellNetravali)

		rout, err := os.Create(cacheSavePath + "/" + id + ".jpg")
		if err != nil {
			log.Debugf("outfileCreate error %s", err)
			return err
		}
		defer rout.Close()

		err = jpeg.Encode(rout, resizedImage, &jpeg.Options{
			Quality: quality,
		})
		if err != nil {
			log.Debugf("imageEncode error %s", err)
			return err
		}

		return ctx.File(cacheSavePath + "/" + id + ".jpg")
	}
}

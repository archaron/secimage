package api

import (
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/chai2010/webp"
	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
	"go.uber.org/zap"
)

type (
	viRequest struct {
		ID            string `params:"id" validate:"required"`
		Size          string `params:"size" validate:"required"`
		File          string `params:"file" validate:"required"`
		width, height uint64
	}

	viParams struct {
		log          *zap.Logger
		cli          *http.Client
		savePath     string
		cachePath    string
		allowedSizes []string
		jpegQuality  int
		webpQuality  float32
		webpLossless bool
		imageType    string
		inputFormat  string
	}
)

var re = regexp.MustCompile(`(\d+)x(\d+)`)

func vi(p viParams) (echo.HandlerFunc, error) {
	allowed := make(map[string]struct{}, len(p.allowedSizes))

	// fill up size hashmap
	for _, size := range p.allowedSizes {
		allowed[size] = struct{}{}
	}

	// creates all needed path, if they're not exists
	for _, item := range []string{p.cachePath, p.savePath} {
		if _, err := os.Stat(item); err != nil {
			if os.IsNotExist(err) {
				p.log.Debug("creating directory",
					zap.String("path", item),
				)
				if err := os.MkdirAll(item, 0777); err != nil {
					p.log.Error("could not create path", zap.String("path", item), zap.Error(err))
					return nil, err
				}
			} else {
				p.log.Error("could not stat path", zap.String("path", item), zap.Error(err))
			}
		}
	}

	return func(ctx echo.Context) error {
		var req viRequest

		var err error
		var img image.Image

		if err := ctx.Bind(&req); err != nil {
			p.log.Error("can't bind/validate request",
				zap.Error(err))

			return ctx.String(http.StatusBadRequest, "FAIL")
		}
		match := re.FindAllStringSubmatch(req.Size, -1)

		if len(match) < 0 || len(match[0]) < 3 {
			p.log.Error("size not match regexp")
			return ctx.String(http.StatusBadRequest, "FAIL SIZE")
		}

		var (
			size   = match[0][0]
			width  = match[0][1]
			height = match[0][2]
		)

		if _, ok := allowed[size]; !ok {
			p.log.Error("size not allowed", zap.String("size", size))
			return ctx.String(http.StatusBadRequest, "Size is not allowed")
		}

		isWebp := strings.HasSuffix(req.File, "webp")

		cacheSavePath := p.cachePath + "/" + size

		var (
			sizedFile    string
			originalFile string
			originalUrl  string
		)

		if isWebp {
			sizedFile = cacheSavePath + "/" + req.ID + ".webp"
		} else {
			sizedFile = cacheSavePath + "/" + req.ID + ".jpg"
		}

		switch p.inputFormat {
		case "webp":
			originalUrl = "https://i.ytimg.com/vi_webp/" + req.ID + "/" + p.imageType + ".webp"
			originalFile = p.savePath + "/" + req.ID + ".webp"
		default:
			originalUrl = "https://i.ytimg.com/vi/" + req.ID + "/" + p.imageType + ".jpg"
			originalFile = p.savePath + "/" + req.ID + ".jpg"
		}

		p.log.Debug("files",
			zap.String("request_file", req.File),
			zap.String("cache_save_path", cacheSavePath),
			zap.String("sized_file", sizedFile),
			zap.String("original_file", originalFile),
			zap.String("original_url", originalUrl),
			zap.String("input_format", p.inputFormat),
			zap.Bool("is_webp", isWebp),
		)

		if _, err = os.Stat(sizedFile); os.IsNotExist(err) {
			p.log.Debug("no sized file")
			req.width, err = strconv.ParseUint(width, 10, 64)
			if err != nil {
				return err
			}

			req.height, err = strconv.ParseUint(height, 10, 64)
			if err != nil {
				return err
			}

			// ensure directory exists
			if _, err := os.Stat(cacheSavePath); err != nil {
				if os.IsNotExist(err) {
					if err := os.MkdirAll(cacheSavePath, 0777); err != nil {
						p.log.Error("can't create folder",
							zap.Error(err))
						return ctx.String(http.StatusBadRequest, "FAIL")
					}
				} else {
					p.log.Error("cannot stat cache save path", zap.Error(err))
					return ctx.String(http.StatusBadRequest, "FAIL")
				}
			}

			if _, err := os.Stat(originalFile); os.IsNotExist(err) {
				p.log.Debug("Downloading",
					zap.String("id", req.ID),
					zap.String("size", size))

				// Get the data
				resp, err := http.Get(originalUrl)
				if err != nil {
					p.log.Debug("could not download file",
						zap.String("file", originalUrl),
						zap.Error(err))
					return err
				}
				defer resp.Body.Close()

				out, err := os.Create(originalFile)
				if err != nil {
					p.log.Debug("could not create file", zap.Error(err))
					return err
				}
				defer out.Close()

				// Write the body to file
				_, err = io.Copy(out, resp.Body)
				if err != nil {
					p.log.Debug("coping file fail",
						zap.Error(err))
					return err
				}

			} else {
				p.log.Debug("cached original",
					zap.String("id", req.ID),
					zap.String("size", size))
			}

			in, err := os.Open(originalFile)
			if err != nil {
				return err
			}
			defer in.Close()

			if p.inputFormat == "webp" {
				img, err = webp.Decode(in)
			} else {
				img, _, err = image.Decode(in)
			}

			// Image decode error catch
			if err != nil {
				p.log.Debug("could not decode image", zap.Error(err))
				return err
			}

			resizedImage := resize.Resize(uint(req.width), uint(req.height), img, resize.MitchellNetravali)

			rout, err := os.Create(sizedFile)
			if err != nil {
				p.log.Debug("could not create store file", zap.Error(err))
				return err
			}
			defer rout.Close()

			if isWebp {
				err = webp.Encode(rout, resizedImage, &webp.Options{
					Lossless: p.webpLossless,
					Quality:  p.webpQuality,
					Exact:    false,
				})
			} else {
				err = jpeg.Encode(rout, resizedImage, &jpeg.Options{
					Quality: p.jpegQuality,
				})
			}

			if err != nil {
				p.log.Debug("could not encode image", zap.Error(err))
				return err
			}

		}

		return ctx.File(sizedFile)
	}, nil
}

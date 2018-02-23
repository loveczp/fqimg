package lib

import (
	"net/http"
	"image"
	"image/jpeg"
	"strconv"
	"image/png"
	"image/gif"
	"golang.org/x/image/bmp"
	"github.com/chai2010/webp"
)

func getQuality(para map[string]string) (int, error) {
	if v, ok := para["q"]; ok {
		return strconv.Atoi(v)
	} else {
		return 0, nil
	}
}

func format_jpeg(resp http.ResponseWriter, req *http.Request, img image.Image, para map[string]string) ( error) {
	quality, err := getQuality(para);
	if err != nil {
		WriteErr(resp,http.StatusBadRequest, err)
	}
	var rgba *image.RGBA
	if nrgba, ok := img.(*image.NRGBA); ok {
		if nrgba.Opaque() {
			rgba = &image.RGBA{
				Pix:    nrgba.Pix,
				Stride: nrgba.Stride,
				Rect:   nrgba.Rect,
			}
		}
	}
	if quality < 1 || 100 < quality {
		quality = 70
	}

	resp.Header().Add("content-type", "image/jpeg")
	if rgba != nil {
		err = jpeg.Encode(resp, rgba, &jpeg.Options{Quality: quality})
	} else {
		err = jpeg.Encode(resp, img, &jpeg.Options{Quality: quality})
	}
	return  nil
}

func format_png(resp http.ResponseWriter, req *http.Request, img image.Image, para map[string]string) (error) {
	resp.Header().Add("content-type", "image/png")
	err := png.Encode(resp, img)
	return err
}

func format_gif(resp http.ResponseWriter, req *http.Request, img image.Image, para map[string]string) (error) {
	quality, err := getQuality(para);
	if err != nil {
		WriteErr(resp,http.StatusBadRequest, err)
	}
	resp.Header().Add("content-type", "image/gif")
	if quality < 1 || 256 < quality {
		quality = 256
	}

	err = gif.Encode(resp, img, &gif.Options{NumColors: quality})
	return nil
}

func format_bmp(resp http.ResponseWriter, req *http.Request, img image.Image, para map[string]string) (error) {
	resp.Header().Add("content-type", "image/bmp")
	err := bmp.Encode(resp, img)
	if err!=nil{
		WriteErr(resp,http.StatusBadRequest, err)
	}
	return nil
}
func format_webp(resp http.ResponseWriter, req *http.Request, img image.Image, para map[string]string) (error) {
	quality, err := getQuality(para);
	if err != nil {
		WriteErr(resp,http.StatusBadRequest, err)
	}
	if quality < 1 || 100 < quality {
		quality = 50
	}
	resp.Header().Add("content-type", "image/webp")
	err = webp.Encode(resp, img, &webp.Options{Lossless: false, Quality:float32(quality)});
	return err
}

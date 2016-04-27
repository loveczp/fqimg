package main

import (
	"github.com/disintegration/imaging"
	"strings"
	"io"
	"image"
	"image/jpeg"
	"image/png"
	"image/gif"
	"golang.org/x/image/bmp"
	"errors"
	"net/http"
	"github.com/chai2010/webp"
	"log"
)

func stringToAnchor(instr  string) imaging.Anchor {
	innerstr := strings.ToLower(instr)
	switch innerstr {
	case "center":
		return imaging.Center
	case "topleft":
		return imaging.TopLeft
	case "top":
		return imaging.Top
	case "topright":
		return imaging.TopRight
	case "left":
		return imaging.Left
	case "right":
		return imaging.Right
	case "bottomleft":
		return imaging.BottomLeft
	case "bottom":
		return imaging.Bottom
	case "bottomright":
		return imaging.BottomRight
	default:
		return imaging.Center
	}
}

func stringToFilter(instr  string) imaging.ResampleFilter {
	innerstr := strings.ToLower(instr)

	switch innerstr {
	case "nearestneighbor":
		return imaging.NearestNeighbor
	case "box":
		return imaging.Box
	case "linear":
		return imaging.Linear
	case "hermite":
		return imaging.Hermite
	case "mitchellnetravali":
		return imaging.MitchellNetravali
	case "catmullrom":
		return imaging.CatmullRom
	case "bspline":
		return imaging.BSpline
	case "gaussian":
		return imaging.Gaussian
	case "bartlett":
		return imaging.Bartlett
	case "lanczos":
		return imaging.Lanczos
	case "hann":
		return imaging.Hann
	case "hamming":
		return imaging.Hamming
	case "blackman":
		return imaging.Blackman
	case "welch":
		return imaging.Welch
	case "cosine":
		return imaging.Cosine
	default:
		return imaging.Lanczos
	}
}

func byte2string(in [16]byte) []byte {
	tmp := make([]byte, 16)
	for _, value := range in {
		tmp = append(tmp, value)
	}
	return tmp[16:]
}


type storage interface {
	storagePut(src io.Reader) (string,error)
	storageGet(key string)  (io.Reader,error)
}



type cache interface {
	cacheGet(url string, desc io.Writer) error
	cachePut(url string, desc io.Reader) error
}



func encode(w http.ResponseWriter, img image.Image, format string,quality int) error {
	//log.Println(format,quality)
	var err error
	switch format {
	case "jpeg":
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
		if  quality < 1 || 100 < quality {
			quality = 70
		}



		w.Header().Add("content-type","image/jpeg")
		if rgba != nil {
			err = jpeg.Encode(w, rgba, &jpeg.Options{Quality: quality})
		} else {
			err = jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
		}

	case "png":
		w.Header().Add("content-type","image/png")
		err = png.Encode(w, img)
	case "gif":
		w.Header().Add("content-type","image/gif")
		if quality < 1 || 256 < quality {
			quality = 256
		}
		err = gif.Encode(w, img, &gif.Options{NumColors: quality})
	case "bmp":
		w.Header().Add("content-type","image/bmp")
		err = bmp.Encode(w, img)
	case "webp":
		w.Header().Add("content-type","image/webp")
		if quality < 1 || 100 < quality {
			quality = 80
		}
		if err = webp.Encode(w, img, &webp.Options{Lossless: false,Quality:float32(quality)}); err != nil {
			log.Fatalln(err)
		}
	default:
		err = errors.New("format not supported")
	}
	return err
}
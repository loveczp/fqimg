package main

import (
	"github.com/disintegration/imaging"
	"strings"
	"io"
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
	storagePut(md5 string, src io.Reader) error
	storageGet(md5 string)  (io.Reader,error)
}



type cache interface {
	cacheGet(url string, desc io.Writer) error
	cachePut(url string, desc io.Reader) error
}




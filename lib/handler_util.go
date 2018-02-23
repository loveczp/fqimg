package lib

import (
	"github.com/disintegration/imaging"
	"strings"
	"errors"
	"strconv"
	"net/http"
	"encoding/json"
	"log"
	"io"
)

func stringToAnchor(instr string) imaging.Anchor {
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

func stringToFilter(instr string) imaging.ResampleFilter {
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

func checkStrength(para map[string]string, defaultValue float64) (float64, error) {
	if _, ok := para["s"]; ok {
		if strength, err := strconv.ParseFloat(para["s"], 64); err == nil {
			return strength, nil
		} else {
			return 0, err
		}
	} else {
		return defaultValue, nil;
	}
}

func checkResizeParameter(para map[string]string) (int, int, string, error) {
	intw, wr := strconv.Atoi(para["w"])
	inth, hr := strconv.Atoi(para["h"])
	if (wr != nil || hr != nil || intw <= 0 || inth <= 0) {
		errorstr := `resize/fill/fit  command parameter number or value is not correct
		these command required at least h=xxx  w=xxx two parameters
		and the value of xxx must greater than zero
		`;
		return 0, 0, "", errors.New(errorstr)
	}
	filter_str := ""
	if _, ok := para["f"]; ok {
		filter_str = para["f"]
	}
	return intw, inth, filter_str, nil
}




func WriteErr(resp http.ResponseWriter, status_code int, err error) {
	jsonstr, _ := json.Marshal(map[string]string{"error": err.Error()})
	log.Println(string(jsonstr));
	resp.WriteHeader(status_code)
	io.WriteString(resp, string(jsonstr))
}

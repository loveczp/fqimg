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

func checkStrength(para []string, defaultValue float64) (float64, error) {
	if len(para)>=2{
		if strength, err := strconv.ParseFloat(para[1], 64); err == nil {
			return strength, nil
		} else {
			return 0, err
		}
	} else {
		return defaultValue, nil;
	}
}

func checkResizeParameter(para []string) (int, int, string, error) {
	intw:=0
	inth:=0
	var err error
	var filter string
	if len(para)>=3{
		intw, err= strconv.Atoi(para[1])
		if err!=nil {
			return 0,0,"",err
		}
		inth, err= strconv.Atoi(para[2])

		if err!=nil {
			return 0,0,"",err
		}

		if len(para)>=4{
			filter =para[3]
		}
	}

	if (intw == 0 || inth == 0) {
		errorstr := `resize/fill/fit  command parameter number or value is not correct
		these command required at least h=xxx  w=xxx two parameters
		and the value of xxx must greater than zero
		`;
		return 0, 0, "", errors.New(errorstr)
	}
	return intw, inth, filter, nil
}




func WriteErr(resp http.ResponseWriter, status_code int, err error) {
	jsonstr, _ := json.Marshal(map[string]string{"error": err.Error()})
	log.Println(string(jsonstr));
	resp.WriteHeader(status_code)
	io.WriteString(resp, string(jsonstr))
}


func parseQueryString(querystr string) [][]string{
	rearr  :=[][]string{}
	if len(querystr)==0{
		return rearr
	}
	paras :=strings.Split(querystr,"&")
	for _,value:=range paras{
		cmd_item := []string{}
		items := strings.Split(value,"=")
		if len(items)>=1{
			cmd_item=append(cmd_item,items[0])
		}
		if len(items)>=2{
			cmd_item=append(cmd_item,items[1])
		}
		rearr=append(rearr,cmd_item)
	}
	return rearr
}
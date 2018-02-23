package lib

import (
	"github.com/disintegration/imaging"
	"strings"
	"net/http"
	"image"
	"io"
	"github.com/deckarep/golang-set"
	"github.com/loveczp/fqimg/store"
	"fmt"
	"net/url"
	"log"
)

var (
	cmds    mapset.Set
	cmd_map = map[string]func(para []string, in image.Image) (image.Image, error){
		"fit":        cmd_fit,
		"fill":       cmd_fill,
		"resize":     cmd_resize,
		"gamma":      cmd_gamma,
		"sigmoid":    cmd_sigmoid,
		"contrast":   cmd_contrast,
		"brightness": cmd_brightness,
		"invert":     cmd_invert,
		"grayscale":  cmd_grayscale,
		"blur":       cmd_blur,
		"sharpen":    cmd_sharpen,
		"rotate90":   cmd_rotate90,
		"rotate180":  cmd_rotate180,
		"rotate270":  cmd_rotate270,
		"flipH":      cmd_flipH,
		"flipV":      cmd_flipV,
		"transpose":  cmd_transpose,
		"mark":       cmd_mark}
	format_map = map[string]func(resp http.ResponseWriter, req *http.Request, img image.Image, para []string) (error){
		"jpeg": format_jpeg,
		"png":  format_png,
		"gif":  format_gif,
		"bmp":  format_bmp,
		"webp": format_webp}
)

func init() {
	cmds = mapset.NewSet()
	for key := range cmd_map {
		cmds.Add(key)
	}
	for key := range format_map {
		cmds.Add(key)
	}
}

func GetHandler(store store.Storage) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		key := req.URL.Path[1:]
		key = strings.TrimPrefix(key, "get/")
		var outImage image.Image
		reader, err := store.Get(key);
		if err != nil {
			WriteErr(resp, http.StatusBadRequest, err)
			return
		}
		outImage, err = imaging.Decode(reader)
		if err != nil {
			WriteErr(resp, http.StatusBadRequest, err)
			return
		}
		ops, format_para, err := getCommands(req)
		if err != nil {
			WriteErr(resp, http.StatusBadRequest, err)
		}

		for _,value:=range(ops) {
			outImage, err = cmd_map[value[0]](value, outImage);
			if err != nil {
				WriteErr(resp, http.StatusBadRequest, err)
				return
			}
		}

		if len(format_para)==0 {
			format_para=append(format_para, "jpeg")
		}
		format_map[format_para[0]](resp, req, outImage, format_para)
		return
	}
}

func getCommands(req *http.Request) ([][]string,[]string , error) {
	raw_querry, err := url.PathUnescape(req.URL.RawQuery)
	if err!=nil{
		return nil,nil,err
	}

	querry_arr:=parseQueryString(raw_querry)

	re_cmds :=[][]string{}
	re_format:=[]string{}
	for _,value:=range(querry_arr) {
		cmd:=[]string{}
		cmd=append(cmd,value[0])
		if len(value)>1{
			paraString := strings.TrimSpace(value[1])
			if (len(paraString) < 1) {
				break
			}
			values:=strings.Split(paraString,"_")
			cmd=append(cmd,values...)
		}

		if _,ok:=cmd_map[value[0]];ok{
			re_cmds=append(re_cmds,cmd)
		}
		if _,ok:=format_map[value[0]];ok{
			re_format=cmd
		}
	}
	return re_cmds, re_format, nil
}

func HelloHandler() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		io.WriteString(resp, "hello")
	}
}

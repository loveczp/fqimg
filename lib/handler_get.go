package lib

import (
	"github.com/disintegration/imaging"
	"strings"
	"net/http"
	"image"
	"io"
	"container/list"
	"github.com/deckarep/golang-set"
	"github.com/loveczp/fqimg/store"
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"log"
)

var (
	cmds    mapset.Set
	cmd_map = map[string]func(para map[string]string, in image.Image) (image.Image, error){
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
	format_map = map[string]func(resp http.ResponseWriter, req *http.Request, img image.Image, para map[string]string) (error){
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
		key = strings.TrimPrefix(key, getAlias+"/")
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
		//log.Println("ops:",ops)
		if err != nil {
			WriteErr(resp, http.StatusBadRequest, err)
		}

		if (ops.Len() == 0) {
			imaging.Encode(resp, outImage, imaging.JPEG)
			return
		}

		for e := ops.Front(); e != nil; e = e.Next() {
			para, _ := e.Value.(map[string]string)
			command := para["c"]
			outImage, err = cmd_map[command](para, outImage);
			if err != nil {
				WriteErr(resp, http.StatusBadRequest, err)
				return
			}
		}

		format_cmd, ok := format_para["c"]
		//fmt.Println("format_cmd",format_cmd)
		if ok == false {
			format_cmd = "jpeg"
			format_para["c"] = "jpeg"
		}
		format_map[format_cmd](resp, req, outImage, format_para)
		return
	}
}

func getCommands(req *http.Request) (*list.List, map[string]string, error) {
	raw_query := req.URL.RawQuery;
	query, err := url.PathUnescape(raw_query)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("the query %s format is wrong", raw_query))
	}
	log.Println("query:", query)
	if len(query) == 0 && len(Conf.DefaultAction) > 2 {
		query = Conf.DefaultAction;
	}
	ops := list.New()
	opts := strings.Split(query, "|")
	formatpara := map[string]string{}
	for i := 0; i < len(opts); i++ {
		paraString := strings.TrimSpace(opts[i])
		if (len(paraString) < 1) {
			break
		}
		paramap := map[string]string{}
		paras := strings.Split(paraString, "&")
		for j := 0; j < len(paras); j++ {
			pairArray := strings.Split(paras[j], "=")

			if (len(paras[j]) < 1 || len(pairArray) != 2) {
				return nil, nil, errors.New(fmt.Sprintf("the parameter %s format is wrong", paras[j]))
			}
			paramap[pairArray[0]] = pairArray[1]
		}
		incom, ok := paramap["c"]
		if (ok && cmds.Contains(incom) == false) {
			return nil, nil, errors.New(fmt.Sprintf("the commond %s is not available", incom))
		}

		if _, ok := cmd_map[incom]; ok {
			ops.PushBack(paramap)
		}

		if _, ok := format_map[incom]; ok {
			formatpara = paramap
		}
	}
	return ops, formatpara, nil
}

func HelloHandler() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		fmt.Println("hello resp")
		io.WriteString(resp, "hello")
	}
}

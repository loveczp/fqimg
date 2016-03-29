package main

import (
	"github.com/disintegration/imaging"
	"strings"
	"strconv"
	"encoding/json"
	"net/http"
	"image"
	"io"
	"container/list"
	"log"
	"errors"
)

func getHandler(resp http.ResponseWriter, req *http.Request) {

	headers := make(map[string]string)
	headers["test"] = "test"

	md5string := strings.Replace(req.URL.Path[1:], "/", "", 100)
	//logrus.Println("request come:", md5string)
	allcommands := map[string]string{"original":"", "fit":"", "fill":"", "resize":"", "gamma":"", "sigmoid":"", "contrast":"", "brightness":"", "invert":"", "grayscale":"", "blur":"", "sharpen":"", "rotate90":"", "rotate180":"", "flipH":"", "flipV":"", "transpose":"", "transverse":""}

	//md5file := "./upload/" + md5string;

	var outImage image.Image
	if reader, err := store.storageGet(md5string); err == nil {
		outImage, err = imaging.Decode(reader)
		if err != nil {
			log.Println("image docode error:" + err.Error());
			io.WriteString(resp, "image docode error:" + err.Error())
			return
		}else {
			//log.Println("outImage docoded:");
		}
	}else {
		jsonstr, _ := json.Marshal(map[string]string{"error": "the image you reqeust does not exist:" + err.Error(), "original":md5string})
		log.Println(string(jsonstr));
		io.WriteString(resp, string(jsonstr))
		return
	}

	ops := list.New()
	//res.WriteHeader()
	opts := strings.Split(req.URL.RawQuery, "|")
	for i := 0; i < len(opts); i++ {
		paraString := strings.TrimSpace(opts[i])

		//if this command string is null then stop all the following command parsing
		if (len(paraString) < 1) {
			break
		}
		paramap := make(map[string]string)
		paras := strings.Split(paraString, "&")
		for j := 0; j < len(paras); j++ {
			pairArray := strings.Split(paras[j], "=")

			//if the parameter format is not foo=bar then  ignore this parameter

			if (len(paras[j]) < 1 || len(pairArray) != 2) {
				jsonstr, _ := json.Marshal(map[string]string{"error": "the parameter format is wrong", "original":paras[j]})

				log.Println(string(jsonstr));
				io.WriteString(resp, string(jsonstr))
				return
			}

			paramap[pairArray[0]] = pairArray[1]
		}

		incom, ok := paramap["c"]
		_, ok2 := allcommands[incom]
		if (ok && ok2 == false) {
			jsonstr, _ := json.Marshal(map[string]string{"error": "the command is not applicable", "original":incom})
			log.Println(string(jsonstr));
			io.WriteString(resp, string(jsonstr))
			return
		}

		ops.PushBack(paramap)
	}

	for e := ops.Front(); e != nil; e = e.Next() {
		v, _ := e.Value.(map[string]string)
		intw, _ := strconv.Atoi(v["w"])
		inth, _ := strconv.Atoi(v["h"])
		filter := v["f"]
		command := v["c"]
		switch command {

		case "original":{
			imaging.Encode(resp, outImage, imaging.JPEG)
			return
		}

		//resize
		case "fit":{
			if error := checkResizeParameter(v); error != nil {
				io.WriteString(resp, error.Error())
				return
			}
			outImage = imaging.Fit(outImage, intw, inth, stringToFilter(filter))
		}
		case "fill":{
			if error := checkResizeParameter(v); error != nil {
				io.WriteString(resp, error.Error())
				return
			}
			outImage = imaging.Fill(outImage, intw, inth, stringToAnchor(v["a"]), stringToFilter(filter))
		}
		case "resize":{
			if error := checkResizeParameter(v); error != nil {
				io.WriteString(resp, error.Error())
				return
			}
			outImage = imaging.Resize(outImage, intw, inth, stringToFilter(filter))
		}


		//adjust
		case "gamma":{
			if value ,err:=checkStrength(v,0.7); err!=nil{
				outImage = imaging.AdjustGamma(outImage, value)
			}else{
				io.WriteString(resp, "gamma strength para error:"+err.Error())
				return

			}
		}
		case "sigmoid":{
			outImage = imaging.AdjustSigmoid(outImage, 0.5, 3.0)
		}

		case "contrast":{
			if value ,err:=checkStrength(v,20); err!=nil{
				io.WriteString(resp, "contrast strength para error:"+err.Error())
				return
			}else{
				outImage = imaging.AdjustContrast(outImage, value)
			}
		}
		case "brightness":{
			if value ,err:=checkStrength(v, 0.5); err!=nil{
				io.WriteString(resp, "brightness strength para error:"+err.Error())
				return
			}else{
				outImage = imaging.AdjustBrightness(outImage,value)
			}

		}
		case "grayscale":{
			outImage = imaging.Grayscale(outImage)
		}
		case "invert":{
			outImage = imaging.Invert(outImage)
		}




		//effects
		case "blur":{

			if value ,err:=checkStrength(v,3.5); err!=nil{
				io.WriteString(resp, "brightness strength para error:"+err.Error())
				return
			}else{
				outImage = imaging.Blur(outImage, value)
			}


		}
		case "sharpen":{
			if value ,err:=checkStrength(v,3.5); err!=nil{
				io.WriteString(resp, "brightness strength para error:"+err.Error())
				return
			}else{
				outImage = imaging.Sharpen(outImage, value)
			}
		}



		//transform
		case "rotate90":{
			outImage = imaging.Rotate90(outImage)
		}
		case "rotate180":{
			outImage = imaging.Rotate180(outImage)
		}
		case "rotate270":{
			outImage = imaging.Rotate270(outImage)
		}
		case "flipH":{
			outImage = imaging.FlipH(outImage)
		}
		case "flipV":{
			outImage = imaging.FlipV(outImage)
		}
		case "transpose":{
			outImage = imaging.Transpose(outImage)
		}
		case "transverse":{
			outImage = imaging.Transverse(outImage)
		}
		}
	}

	//imaging.Encode(res, outImage, imaging.JPEG)
	if (outImage == nil) {
		io.WriteString(resp, "outimage is null")
	}else {
		imaging.Encode(resp, outImage, imaging.JPEG)
	}
	return
}


func checkResizeParameter(para map[string]string) error {
	intw, wr := strconv.Atoi(para["w"])
	inth, hr := strconv.Atoi(para["h"])
	if (wr != nil || hr != nil || intw <= 0 || inth <= 0) {
		errorstr := `resize/fill/fit  command parameter number or value is not correct
		these command required at least h=xxx  w=xxx two parameters
		and the value of xxx must greater than zero
		`;
		return errors.New(errorstr)
	}
	return nil
}

func checkStrength(para map[string]string,defaultValue float64) (float64,error) {
	if _, ok :=para["s"]; ok {
		if strength, err := strconv.ParseFloat(para["s"],64);err ==nil{
			return strength,nil
		}else {
			return 0,err
		}
	}else {
		return defaultValue,nil;
	}
}
package imageserverlib

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
	"time"
	"github.com/deckarep/golang-set"
	"image/draw"
	"image/color"
	"os"
	"github.com/hashicorp/golang-lru"
)

var (
	cmdstr = []interface{}{"ori", "fit", "fill", "resize", "gamma", "sigmoid", "contrast",
		"brightness", "invert", "grayscale", "blur", "sharpen", "rotate90", "rotate180",
		"rotate270", "flipH", "flipV", "transpose", "jpeg", "png", "gif", "bmp", "webp", "mark"}
	cmds = mapset.NewSetFromSlice(cmdstr)

	formats = mapset.NewSetFromSlice([]interface{}{"jpeg", "png", "gif", "bmp", "webp"})
	offps = mapset.NewSetFromSlice([]interface{}{"lu", "ru", "ld", "rd"})

	fileCache *(lru.Cache)
)

type imageCacheItem struct {
	key         string
	filePath    string
	contentType string
}

func removeFile(key interface{}, value interface{})  {
	citem := value.(imageCacheItem)
	log.Println("romve cache item go routin out :",citem.key);
	go func(filePath string) {
		log.Println("romve cache item go routin in :",filePath);
		if _ ,err:=os.Stat(filePath);os.IsExist(err) {
			err:=os.Remove(filePath)
			log.Panic("remove cache item error:",err.Error())
		}
	}(citem.filePath)
}

func GetHandler(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	//log.Println(headers);
	if Conf.Headers != nil && len(Conf.Headers) > 0 {
		for key, value := range Conf.Headers {
			resp.Header().Add(key, value)
		}
	}

	/**	get data from local file cache start	 */
	log.Println("get cache, file cache key :", req.URL.String());

	if fileCache==nil{
		var err error;
		fileCache, err = lru.NewWithEvict(Conf.FileCacheSize, removeFile)
		if err !=nil{
			log.Panic("cache create error :",err)
		}
	}
	log.Println("fileCache keys :", (*fileCache).Keys());
	log.Println("fileCache keys Number:", (*fileCache).Len());

	if item, ok := (*fileCache).Get(req.URL.String()); ok {
		citem := item.(imageCacheItem)
		if cfile, err := os.Open(citem.filePath); err == nil {
			resp.Header().Add("content-type", citem.contentType)
			io.Copy(resp, cfile);
			log.Println("data from file cache:", req.URL.String());
			defer  cfile.Close()
			return;
		} else {
			//when cache error, get data from storage
			(*fileCache).Remove(req.URL.String())
			//io.WriteString(resp, "get cache data error:" + err.Error())
			log.Println("data error from file cache:", req.URL.String());
		}
	}
	/**	get data from local file cache end	 */

	log.Println("data miss from cache:", req.URL.String());

	key := req.URL.Path[1:]

	//log.Println("handler_get   key:", key)

	//md5file := "./upload/" + md5string;

	var outImage image.Image
	reader, err := store.storageGet(key);
	if err != nil {
		jsonstr, _ := json.Marshal(map[string]string{"error": "the image you reqeust does not exist:" + err.Error(), "original":key})
		log.Println(string(jsonstr));
		io.WriteString(resp, string(jsonstr))
		return
	}
	outImage, err = imaging.Decode(reader)

	if err != nil {
		log.Println("image docode error:" + err.Error());
		io.WriteString(resp, "image docode error:" + err.Error())
		return
	}



	//if no action presented , add default action
	query := req.URL.RawQuery;
	if len(query) == 0 && len(Conf.DefaultAction) > 2 {
		query = Conf.DefaultAction;
		//log.Println("defaultAction:"+default_action)
	}

	//build  commands
	ops := list.New()
	opts := strings.Split(query, "|")
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
		//_, ok2 := allcommands[incom]
		if (ok && cmds.Contains(incom) == false) {
			jsonstr, _ := json.Marshal(map[string]string{"error": "the command is not applicable", "original":incom})
			log.Println(string(jsonstr));
			io.WriteString(resp, string(jsonstr))
			return
		}

		ops.PushBack(paramap)
	}

	if (ops.Len() == 0) {
		imaging.Encode(resp, outImage, imaging.JPEG)
		return
	}

	for e := ops.Front(); e != nil; e = e.Next() {
		v, _ := e.Value.(map[string]string)
		intw, _ := strconv.Atoi(v["w"])
		inth, _ := strconv.Atoi(v["h"])
		filter := v["f"]
		command := v["c"]
		switch command {

		//resize
		case "ori":{
			imaging.Encode(resp, outImage, imaging.JPEG)
		}
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
			if value, err := checkStrength(v, 0.7); err != nil {
				io.WriteString(resp, "gamma strength para error:" + err.Error())
				return
			} else {
				outImage = imaging.AdjustGamma(outImage, value)
			}
		}
		case "sigmoid":{
			outImage = imaging.AdjustSigmoid(outImage, 0.5, 3.0)
		}

		case "contrast":{
			if value, err := checkStrength(v, 20); err != nil {
				io.WriteString(resp, "contrast strength para error:" + err.Error())
				return
			} else {
				outImage = imaging.AdjustContrast(outImage, value)
			}
		}
		case "brightness":{
			if value, err := checkStrength(v, 0.5); err != nil {
				io.WriteString(resp, "brightness strength para error:" + err.Error())
				return
			} else {
				outImage = imaging.AdjustBrightness(outImage, value)
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

			if value, err := checkStrength(v, 3.5); err != nil {
				io.WriteString(resp, "brightness strength para error:" + err.Error())
				return
			} else {
				outImage = imaging.Blur(outImage, value)
			}
		}
		case "sharpen":{
			if value, err := checkStrength(v, 3.5); err != nil {
				io.WriteString(resp, "brightness strength para error:" + err.Error())
				return
			} else {
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
		case "mark":{
			if outImage, err = mark(v, outImage); err != nil {
				io.WriteString(resp, err.Error())
				return;
			}
		}
		}
	}

	var quality int;
	var command string;
	var ok bool = false
	for e := ops.Front(); e != nil; e = e.Next() {
		v, _ := e.Value.(map[string]string)
		quality, err = strconv.Atoi(v["q"])
		command, ok = v["c"]
		//_, tok := formats[command]

		if ok && formats.Contains(command) {
			break
		} else {
			command = "jpeg"
		}
	}
	//log.Println(outImage.Bounds().String())

	if (outImage == nil) {
		io.WriteString(resp, "outimage is null")
	} else {
		encode(resp, req, outImage, command, quality);
	}

	elapsed := time.Since(start)

	accessLog.Printf("%-10s  %-20s %-50s", elapsed, req.RemoteAddr, req.URL)
	return
}

func mark(para map[string]string, in image.Image) (image.Image, error) {

	if _, ok := para["mid"]; ok {
		if waterMarker, ok := markHash[para["mid"]]; ok {
			offx := 10
			offy := 10
			offxTemp, okx := para["offx"]
			offyTemp, oky := para["offy"]
			if (okx && oky) {
				if _, err := strconv.Atoi(offxTemp); err == nil {
					offx, _ = strconv.Atoi(offxTemp);
				} else {
					log.Println(err)
				}

				if _, err := strconv.Atoi(offyTemp); err == nil {
					offy, _ = strconv.Atoi(offyTemp);
				} else {
					log.Println(err)
				}
			}

			offp := "rd"
			offpTemp, okp := para["offp"]
			if (okp && offps.Contains(offpTemp)) {
				offp = offpTemp
			}

			var bound = image.ZR
			switch offp {
			case "lu":{
				bound = waterMarker.Bounds().Add(image.Pt(offx, offy))
			}
			case "rd":{
				max := in.Bounds().Sub(image.Pt(offx, offy)).Max
				min := image.Pt(max.X - waterMarker.Bounds().Dx(), max.Y - waterMarker.Bounds().Dy())
				bound.Max = max;
				bound.Min = min;

			}
			case "ld":{
				max := image.Pt(waterMarker.Bounds().Max.X + offx, in.Bounds().Max.Y - offy)
				min := image.Pt(waterMarker.Bounds().Min.X + offx, in.Bounds().Max.Y - offy - waterMarker.Bounds().Dy())
				bound.Max = max;
				bound.Min = min;
			}
			case "ru":{
				max := image.Pt(in.Bounds().Max.X - offx, waterMarker.Bounds().Max.Y + offy)
				min := image.Pt(in.Bounds().Max.X - offx - waterMarker.Bounds().Dx(), waterMarker.Bounds().Min.Y + offy)
				bound.Max = max;
				bound.Min = min;
			}
			}





			//offset := image.Pt(400, 400)
			log.Println(waterMarker.Bounds())
			log.Println(bound)
			log.Println(waterMarker.Bounds().Add(image.Pt(offx, offy)))

			//waterMarkerPost:=image.NewRGBA(waterMarker.Bounds()).Opaque()
			var alpha uint8 = 100;

			if alphaTemp, ok := para["alpha"]; ok {
				if alphaTempInt, err := strconv.ParseUint(alphaTemp, 10, 8); err == nil {
					alpha = uint8(alphaTempInt)
				}
			}

			mask := image.NewUniform(color.Alpha{alpha})
			b := in.Bounds()
			m := image.NewRGBA(b)
			draw.Draw(m, b, in, image.ZP, draw.Src)
			//draw.Draw(m, waterMarker.Bounds().Add(image.Pt(offx,offy)), waterMarker, image.ZP, draw.Over)
			draw.DrawMask(m, bound, waterMarker, image.ZP, mask, image.ZP, draw.Over)
			return m, nil
		} else {
			return nil, errors.New("marker image id is null");
		}
	} else {
		return nil, errors.New("marker image id is null");
	}
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
package imageserverlib

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
	"net"
	"strconv"
	"os"
	"encoding/base64"
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
	storagePut(src io.Reader) (string, error)
	storageGet(key string) (io.Reader, error)
}

type cache interface {
	cacheGet(url string, desc io.Writer) error
	cachePut(url string, desc io.Reader) error
}

func encode(w http.ResponseWriter, req *http.Request, img image.Image, format string, quality int) error {
	//log.Println(format,quality)


	//put to lru cache   start
	cPath := Conf.FileCacheDir + base64.StdEncoding.EncodeToString([]byte(req.URL.String()));
	var tempFile *(os.File)
	var err error
	if _, err := os.Stat(cPath); !os.IsExist(err) {
		tempFile, err = os.Create(cPath);
		defer (*tempFile).Close()
	}

	log.Println("cache set data , key ", req.URL.String())
	citem := imageCacheItem{filePath:cPath, key:req.URL.String()}

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
		if quality < 1 || 100 < quality {
			quality = 70
		}

		w.Header().Add("content-type", "image/jpeg")
		citem.contentType = "image/jpeg"
		if rgba != nil {
			err = jpeg.Encode(w, rgba, &jpeg.Options{Quality: quality})
			if tempFile != nil {
				err = jpeg.Encode(tempFile, rgba, &jpeg.Options{Quality: quality})
			}
		} else {
			err = jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
			if tempFile != nil {
				err = jpeg.Encode(tempFile, img, &jpeg.Options{Quality: quality})
			}
		}

	case "png":
		w.Header().Add("content-type", "image/png")
		citem.contentType = "image/png"

		err = png.Encode(w, img)
		if tempFile != nil {
			err = png.Encode(tempFile, img)
		}
	case "gif":
		w.Header().Add("content-type", "image/gif")
		citem.contentType = "image/gif"
		if quality < 1 || 256 < quality {
			quality = 256
		}
		err = gif.Encode(w, img, &gif.Options{NumColors: quality})
		if tempFile != nil {
			err = gif.Encode(tempFile, img, &gif.Options{NumColors: quality})
		}
	case "bmp":
		w.Header().Add("content-type", "image/bmp")
		citem.contentType = "image/bmp"
		err = bmp.Encode(w, img)
		if tempFile != nil {
			err = bmp.Encode(tempFile, img)
		}
	case "webp":
		w.Header().Add("content-type", "image/webp")
		citem.contentType = "image/webp"
		if quality < 1 || 100 < quality {
			quality = 50
		}
		err = webp.Encode(w, img, &webp.Options{Lossless: false, Quality:float32(quality)});
		if tempFile != nil {
			err = webp.Encode(tempFile, img, &webp.Options{Lossless: false, Quality:float32(quality)});
		}
	default:
		err = errors.New("format not supported")
	}

	(*fileCache).Add(req.URL.String(), citem)
	return err
}

func parseIp(ips []string) []interface{} {
	re := make([]interface{}, len(ips));
	for i := 0; i < len(ips); i++ {
		if strings.Contains(ips[i], "/") {
			parts := strings.Split(ips[i], "/");
			netipTemp := parts[0];
			masklenStrTemp := parts[1];
			netip := net.ParseIP(netipTemp);
			maskLenTemp, err1 := strconv.Atoi(masklenStrTemp);
			if err1 != nil {
				log.Fatalln(err1);
				return nil;
			}
			mask := net.CIDRMask(maskLenTemp, 32)
			re[i] = net.IPNet{netip, mask}
		} else {
			re[i] = net.ParseIP(ips[i]);
		}
	}
	return re
}

func ipPass(req *http.Request) bool {
	adds := strings.Split(req.RemoteAddr, ":")
	ip := net.ParseIP(adds[0]);
	//log.Print("req remote ip :", ip)
	if len(Conf.UploadAllowedInterface) > 0 {

		for i := 0; i < len(Conf.UploadAllowedInterface); i++ {
			switch Conf.UploadAllowedInterface[i].(type) {
			case net.IPNet:
				ipnet := Conf.UploadAllowedInterface[i].(net.IPNet)
				if ipnet.Contains(ip) {
					//log.Println(ip.String()," match allow: ",ipnet.String())
					return true;
				}
			case net.IP:
				thisip := Conf.UploadAllowedInterface[i].(net.IP)
				if thisip.Equal(ip) {
					//log.Println(ip.String()," match allow: ",thisip.String())
					return true;
				}
			}
		}
		//log.Println(ip.String()," miss all allow: ",conf.UploadAllowedInterface)
		return false;
	}

	if len(Conf.UploadDenyInterface) > 0 {
		for i := 0; i < len(Conf.UploadDenyInterface); i++ {
			switch Conf.UploadDenyInterface[i].(type) {
			case net.IPNet:
				ipnet := Conf.UploadDenyInterface[i].(net.IPNet)
				if ipnet.Contains(ip) {
					//log.Println(ip.String()," match deny: ",ipnet.String())
					return false;
				}
			case net.IP:
				thisip := Conf.UploadDenyInterface[i].(net.IP)
				if thisip.Equal(ip) {
					//log.Println(ip.String()," match deny: ",thisip.String())
					return false;
				}
			}
		}
		//log.Println(ip.String()," miss all deny: ",conf.UploadDenyInterface)
		return true;
	}

	//log.Println(ip.String()," miss all ip filter")
	return true;
}

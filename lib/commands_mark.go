package lib

import (
	"image"
	"log"
	"image/color"
	"strconv"
	"image/draw"
	"errors"
	"github.com/deckarep/golang-set"
)

var offps = mapset.NewSetFromSlice([]interface{}{"lu", "ru", "ld", "rd"})

func cmd_mark(para map[string]string, in image.Image) (image.Image, error) {

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
			case "lu":
				{
					bound = waterMarker.Bounds().Add(image.Pt(offx, offy))
				}
			case "rd":
				{
					max := in.Bounds().Sub(image.Pt(offx, offy)).Max
					min := image.Pt(max.X-waterMarker.Bounds().Dx(), max.Y-waterMarker.Bounds().Dy())
					bound.Max = max;
					bound.Min = min;
				}
			case "ld":
				{
					max := image.Pt(waterMarker.Bounds().Max.X+offx, in.Bounds().Max.Y-offy)
					min := image.Pt(waterMarker.Bounds().Min.X+offx, in.Bounds().Max.Y-offy-waterMarker.Bounds().Dy())
					bound.Max = max;
					bound.Min = min;
				}
			case "ru":
				{
					max := image.Pt(in.Bounds().Max.X-offx, waterMarker.Bounds().Max.Y+offy)
					min := image.Pt(in.Bounds().Max.X-offx-waterMarker.Bounds().Dx(), waterMarker.Bounds().Min.Y+offy)
					bound.Max = max;
					bound.Min = min;
				}
			}

			log.Println(waterMarker.Bounds())
			log.Println(bound)
			log.Println(waterMarker.Bounds().Add(image.Pt(offx, offy)))
			var alpha uint8 = 10;
			if alphaTemp, ok := para["alpha"]; ok {
				if alphaTempInt, err := strconv.ParseUint(alphaTemp, 10, 8); err == nil {
					alpha = uint8(alphaTempInt)
				}
			}

			mask := image.NewUniform(color.Alpha{alpha})
			b := in.Bounds()
			m := image.NewRGBA(b)
			draw.Draw(m, b, in, image.ZP, draw.Src)
			draw.DrawMask(m, bound, waterMarker, image.ZP, mask, image.ZP, draw.Over)
			return m, nil
		} else {
			return nil, errors.New("marker image id is null");
		}
	} else {
		return nil, errors.New("marker image id is null");
	}
}

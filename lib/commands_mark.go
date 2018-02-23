package lib

import (
	"image"
	"image/color"
	"strconv"
	"image/draw"
	"errors"
	"github.com/deckarep/golang-set"
)

var offps = mapset.NewSetFromSlice([]interface{}{"lu", "ru", "ld", "rd"})

func cmd_mark(para []string, in image.Image) (image.Image, error) {
	mid:=""
	offx := 10
	offy := 10
	offp := "rd"
	alpha:= 255
	var err error
	if len(para)<2{
		return nil, errors.New("mark parameter error");
	}else{
		mid=para[1]
		if len(para)>=4{
			offxTemp := para[2]
			offx, err = strconv.Atoi(offxTemp)
			if err !=nil {
				return nil, errors.New("marker parameter offx err:"+offxTemp);
			}
			offyTemp := para[3]
			offy, err = strconv.Atoi(offyTemp)
			if err !=nil {
				return nil, errors.New("marker parameter offy err:"+offyTemp);
			}
		}

		if len(para)>5{
			offp=para[4]
		}

		if len(para)>6{
			alpha, err = strconv.Atoi(para[5])
			if err!=nil {
				return nil, errors.New("alpha value invalid:"+para[5]);
			}

			if alpha>255 || alpha<1{
				return nil, errors.New("alpha value out of range :"+para[5]+" should be between 0 and 256");
			}
		}
	}

	waterMarker, ok := markHash[mid]
	if ok==false{
		return nil, errors.New("the marker does not exsit :"+mid);
	}
	if offps.Contains(offp)==false {
		return nil, errors.New("the offp does not availble :"+mid+",should be one of 'lu ru ld rd' ");
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


	mask := image.NewUniform(color.Alpha{uint8(alpha)})
	b := in.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, in, image.ZP, draw.Src)
	draw.DrawMask(m, bound, waterMarker, image.ZP, mask, image.ZP, draw.Over)
	return m, nil

}

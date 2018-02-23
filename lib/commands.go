package lib

import (
	"image"
	"github.com/disintegration/imaging"
)

func cmd_fit(para map[string]string, in image.Image) (image.Image, error) {
	if intw, inth, filter_str, err := checkResizeParameter(para); err != nil {
		return  nil,err
	} else {
		out:=imaging.Fit(in, intw, inth, stringToFilter(filter_str))
		return  out,nil
	}
}

func cmd_fill(para map[string]string, in image.Image) (image.Image, error) {

	if intw, inth, filter_str, err := checkResizeParameter(para); err != nil {
		return  nil,err
	} else {
		anchar_str := ""
		if _, ok := para["a"]; ok {
			anchar_str = para["a"]
		}
		out:=imaging.Fill(in, intw, inth, stringToAnchor(anchar_str), stringToFilter(filter_str))
		return  out,nil
	}
}

func cmd_resize(para map[string]string, in image.Image) (image.Image, error) {
	if intw, inth, filter_str, err := checkResizeParameter(para); err != nil {
		return  nil,err
	} else {
		out:=imaging.Resize(in, intw, inth, stringToFilter(filter_str))
		return  out,nil
	}
}

func cmd_gamma(para map[string]string, in image.Image) (image.Image, error) {
	if value, err := checkStrength(para, 0.7); err != nil {
		return  nil,err
	} else {
		out:=imaging.AdjustGamma(in, value)
		return out, nil
	}

}

func cmd_sigmoid(para map[string]string, in image.Image) (image.Image, error) {
	out:=imaging.AdjustSigmoid(in, 0.5, 3.0)
	return out,  nil
}

func cmd_contrast(para map[string]string, in image.Image) (image.Image, error) {
	if value, err := checkStrength(para, 20); err != nil {
		return nil,err
	} else {
		out:=imaging.AdjustContrast(in, value)
		return  out, nil
	}
}

func cmd_brightness(para map[string]string, in image.Image) (image.Image, error) {
	if value, err := checkStrength(para, 0.5); err != nil {
		return nil,err
	} else {
		out:=imaging.AdjustBrightness(in, value)
		return out,nil
	}
}

func cmd_invert(para map[string]string, in image.Image) (image.Image, error) {
	out:=imaging.Invert(in)
	return out, nil
}

func cmd_grayscale(para map[string]string, in image.Image) (image.Image, error) {
	out:=imaging.Grayscale(in)
	return out, nil
}

func cmd_blur(para map[string]string, in image.Image) (image.Image, error) {
	if value, err := checkStrength(para, 3.5); err != nil {
		return nil, err
	} else {
		out:=imaging.Blur(in, value)
		return  out, nil
	}
}

func cmd_sharpen(para map[string]string, in image.Image) (image.Image, error) {
	if value, err := checkStrength(para, 3.5); err != nil {
		return nil,err
	} else {
		out:=imaging.Sharpen(in, value)
		return out, nil
	}
}

func cmd_rotate90(para map[string]string, in image.Image) (image.Image, error) {
	out:=imaging.Rotate90(in)
	return out, nil
}

func cmd_rotate180(para map[string]string, in image.Image) (image.Image, error) {
	out:=imaging.Rotate180(in)
	return out, nil
}

func cmd_rotate270(para map[string]string, in image.Image) (image.Image, error) {
	out:=imaging.Rotate270(in)
	return out, nil
}

func cmd_flipH(para map[string]string, in image.Image) (image.Image, error) {
	out:=imaging.FlipH(in)
	return out, nil
}

func cmd_flipV(para map[string]string, in image.Image) (image.Image, error) {
	out:=imaging.FlipV(in)
	return out, nil
}

func cmd_transpose(para map[string]string, in image.Image) (image.Image, error) {
	out:=imaging.Transpose(in)
	return out, nil
}

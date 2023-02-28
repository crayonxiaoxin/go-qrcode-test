package utils

import (
	"errors"
	"image"
	"os"

	"golang.org/x/image/draw"
)

func OpenImg(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err2 := image.Decode(f)
	if err2 != nil {
		return nil, err2
	}
	return img, nil
}

// size - 宽高都一样；size,size - 宽,高
func ResizeImg(img image.Image, size ...int) (image.Image, error) {
	sizeLen := len(size)
	var width int
	var height int
	if sizeLen == 0 {
		return nil, errors.New("size must not be 0")
	} else if sizeLen == 1 {
		width = size[0]
		height = width
	} else {
		width = size[0]
		height = size[1]
	}
	srcWidth := img.Bounds().Dx()
	srcHeight := img.Bounds().Dy()
	var dstWidth int
	var dstHeight int
	if srcHeight > srcWidth {
		dstHeight = height
		dstWidth = dstHeight * (srcWidth / srcHeight)
	} else {
		dstWidth = width
		dstHeight = dstWidth / (srcWidth / srcHeight)
	}
	dst := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst, nil
}

package main

import (
	"bytes"
	"fmt"
	"gopkg.in/gographics/imagick.v2/imagick"
	"image"
	"image/jpeg"
	"image/png"
	_ "io/ioutil"
	"os"
	"strings"
)

func main() {
	file := os.Args[1]
	infile, err := os.Open(file)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer infile.Close()

	src, format, err := image.Decode(infile)
	if err != nil {
		fmt.Println(err.Error())
	}

	//gray scale optimization
	grayScale(src, file, format)
	magic(src, file, format)

}

func magic(img image.Image, file string, format string) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	buf := new(bytes.Buffer)
	if format == "png" {
		png.Encode(buf, img)
	} else if format == "jpeg" {
		jpeg.Encode(buf, img, &jpeg.Options{jpeg.DefaultQuality})
	}

	err := mw.ReadImageBlob(buf.Bytes())
	if err != nil {
		fmt.Println(err.Error())
	}

	brightContrastFile := strings.Join(strings.Split(file, "."), "_bright_contrast_magic.")
	mw.DespeckleImage()
	mw.BrightnessContrastImage(20, 100)
	mw.WriteImage(brightContrastFile)

	mw.SetImageType(imagick.IMAGE_TYPE_GRAYSCALE)

	grayFile := strings.Join(strings.Split(file, "."), "_gray_magic.")

	mw.WriteImage(grayFile)

}

func grayScale(img image.Image, file string, format string) {

	bounds := img.Bounds()
	gray := image.NewGray(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			oldColor := img.At(x, y)
			grayColor := gray.ColorModel().Convert(oldColor)
			gray.Set(x, y, grayColor)
		}
	}

	grayFile := strings.Join(strings.Split(file, "."), "_gray.")

	outfile, err := os.Create(grayFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer outfile.Close()

	if format == "png" {
		png.Encode(outfile, gray)
	} else if format == "jpeg" {
		jpeg.Encode(outfile, gray, &jpeg.Options{jpeg.DefaultQuality})
	}
}

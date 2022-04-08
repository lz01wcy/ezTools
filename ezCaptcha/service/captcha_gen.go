package main

import (
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"math/rand"
	"strconv"
)

func getLowBrightnessRGBA() color.RGBA {
	r, g, b := color.YCbCrToRGB(uint8(rand.Intn(128)), uint8(rand.Intn(256)), uint8(rand.Intn(256)))
	return color.RGBA{R: r, G: g, B: b, A: 255}
}
func getImg() gocv.Mat {
	mat := <-imageChannel
	background.CopyTo(&mat)
	return mat
}
func genFormula() (int, int, *operator) {
	opIndex := rand.Intn(4)
	op := dic[opIndex]
	switch opIndex {
	// for + -
	case 0, 1:
		doubleDigits := rand.Intn(90) + 10
		singleDigit := rand.Intn(8) + 2
		return doubleDigits, singleDigit, op
	// for *
	case 2:
		doubleDigits := rand.Intn(10) + 10
		singleDigit := rand.Intn(1) + 2
		return doubleDigits, singleDigit, op
	// for /
	case 3:
		doubleDigits := rand.Intn(10) + 10
		singleDigit := rand.Intn(1) + 2
		return doubleDigits * singleDigit, doubleDigits, op
	}
	// never happened
	return 0, 0, nil
}

// Bottom-left corner of the text string in the image.https://docs.opencv.org/4.x/d6/d6e/group__imgproc__draw.html#ga5126f47f883d730f633d74f07456c576
func genLocation(xPadding, maxX, maxY int, str string) image.Point {
	size := gocv.GetTextSize(str, gocv.FontHersheyScriptComplex, 1, 4)
	var x int
	if maxX-size.X <= 0 {
		x = xPadding
	} else {
		x = rand.Intn(maxX-size.X) + xPadding
	}
	y := rand.Intn(maxY-size.Y) + size.Y
	return image.Point{
		X: x,
		Y: y,
	}
}

func gen() ([]byte, string, error) {
	a, b, op := genFormula()
	mat := getImg()
	for i := 0; i < 10; i++ {
		gocv.Line(&mat,
			image.Point{X: rand.Intn(width / 2), Y: rand.Intn(height)},
			image.Point{X: rand.Intn(width/2) + width/2, Y: rand.Intn(height)},
			getLowBrightnessRGBA(),
			rand.Intn(3)+1)
	}
	matA := getImg()
	gocv.Blur(mat, &matA, image.Point{X: 8, Y: 8})
	imageChannel <- mat
	//* 20,23 //高度是23
	gocv.PutText(&matA, strconv.Itoa(a),
		genLocation(0, 80, height, strconv.Itoa(a)), gocv.FontHersheySimplex, 1, getLowBrightnessRGBA(), 4)
	gocv.PutText(&matA, op.o,
		genLocation(80, 30, height, op.o), gocv.FontHersheySimplex, 1, getLowBrightnessRGBA(), 4)
	gocv.PutText(&matA, strconv.Itoa(b),
		genLocation(110, 30, height, strconv.Itoa(b)), gocv.FontHersheySimplex, 1, getLowBrightnessRGBA(), 4)
	gocv.PutText(&matA, "=",
		genLocation(140, 30, height, "="), gocv.FontHersheySimplex, 1, getLowBrightnessRGBA(), 4)
	gocv.PutText(&matA, "?",
		genLocation(170, 30, height, "?"), gocv.FontHersheySimplex, 1, getLowBrightnessRGBA(), 4)
	buffer, err := gocv.IMEncode(gocv.PNGFileExt, matA)
	imageChannel <- matA
	if err != nil {
		return nil, "", err
	}
	return buffer.GetBytes(), strconv.Itoa(op.f(a, b)), nil
}

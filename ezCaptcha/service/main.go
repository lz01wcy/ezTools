package main

import (
	"fmt"
	"github.com/Anveena/ezTools/ezConfig"
	"github.com/Anveena/ezTools/ezPasswordEncoder"
	"gocv.io/x/gocv"
)

type ezCaptchaConfig struct {
	MaxImageCommon int
	GRPCPort       int
	CrtPathEncode  string
	KeyPathEncode  string
	CrtPath        string
	KeyPath        string
}

func (c *ezCaptchaConfig) Check() error {
	if c.CrtPathEncode == "" {
		return fmt.Errorf("CrtPathEncode没配置")
	}
	password, e := ezPasswordEncoder.GetPasswordFromEncodedStr(c.CrtPathEncode)
	if e != nil {
		return e
	}
	c.CrtPath = password
	if c.KeyPathEncode == "" {
		return fmt.Errorf("KeyPathEncode没配置")
	}
	password, e = ezPasswordEncoder.GetPasswordFromEncodedStr(c.KeyPathEncode)
	if e != nil {
		return e
	}
	c.KeyPath = password

	data := make([]byte, width*height*4)
	for i := 0; i < len(data); i++ {
		data[i] = 255
	}
	mat, err := gocv.NewMatFromBytes(height, width, gocv.MatTypeCV8UC4, data)
	if err != nil {
		panic(err)
	}
	background = mat
	imageChannel = make(chan gocv.Mat, c.MaxImageCommon)
	for i := 0; i < c.MaxImageCommon; i++ {
		mat = gocv.NewMatWithSize(height, width, gocv.MatTypeCV8UC4)
		background.CopyTo(&mat)
		imageChannel <- mat
	}
	return nil
}

var captchaConfig = &ezCaptchaConfig{}

type operator struct {
	o string
	f func(a, b int) int
}

var addition = &operator{
	o: "+",
	f: func(a, b int) int {
		return a + b
	},
}
var subtraction = &operator{
	o: "-",
	f: func(a, b int) int {
		return a - b
	},
}
var multiplication = &operator{
	o: "*",
	f: func(a, b int) int {
		return a * b
	},
}
var division = &operator{
	o: "/",
	f: func(a, b int) int {
		return a / b
	},
}

var dic = []*operator{
	addition, subtraction, multiplication, division,
}
var background gocv.Mat

const (
	height = 40
	width  = 200
)

//防止被爆破炸内存.
var imageChannel chan gocv.Mat

func main() {
	if err := ezConfig.ReadConf(captchaConfig); err != nil {
		panic(err)
	}
	if err := captchaConfig.Check(); err != nil {
		panic(err)
	}
	startGRPCService()
}

package main

import (
	"fmt"
	"github.com/Anveena/ezTools/config"
	"github.com/Anveena/ezTools/password"
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

func (c *ezCaptchaConfig) Check() {
	if c.CrtPathEncode == "" {
		panic(fmt.Errorf("CrtPathEncode没配置"))
	}
	pwd, e := password.Decode(c.CrtPathEncode)
	if e != nil {
		panic(e)
	}
	c.CrtPath = pwd
	if c.KeyPathEncode == "" {
		panic(fmt.Errorf("KeyPathEncode没配置"))
	}
	pwd, e = password.Decode(c.KeyPathEncode)
	if e != nil {
		panic(e)
	}
	c.KeyPath = pwd

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

// 防止被爆破炸内存.
var imageChannel chan gocv.Mat

func main() {
	config.ReadConf(captchaConfig)
	startGRPCService()
}

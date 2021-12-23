package main

import (
	"github.com/Anveena/ezTools/ezLog"
	"github.com/Anveena/ezTools/ezPasswordEncoder"
	"time"
)

func main() {
	a, _ := ezPasswordEncoder.EncodePassword("127.0.0.1:12345")
	ezLog.SetUpEnv(&ezLog.EZLoggerModel{
		LogLevel: 0,
		AppName:  "大番薯",
		DingTalkModel: struct {
			Enable                 bool
			SecretKeyEncodedString string
			URLEncodedString       string
			Mobiles                []string
		}{
			Enable: false,
		},
		GRPCModel: struct {
			ClientCounts     int
			URLEncodedString string
		}{
			1,
			a,
		},
	})
	i := uint64(0)
	for {
		s := time.Now()
		for j := 0; j < 10000; j++ {
			i++
			ezLog.D("hello", "world")
			ezLog.DWithTag("hello", "world", i)
		}
		println(time.Now().Sub(s).String())
		//time.Sleep(12345*time.Millisecond)
	}
}

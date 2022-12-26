package main

import (
	"github.com/Anveena/ezTools/log"
	"github.com/Anveena/ezTools/password"
	"time"
)

func main() {
	a, _ := password.Encode("127.0.0.1:12345")
	log.SetUpEnv(&log.EZLoggerModel{
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
			log.D("hello", "world")
			log.DWithTag("hello", "world", i)
		}
		println(time.Now().Sub(s).String())
		//time.Sleep(12345*time.Millisecond)
	}
}

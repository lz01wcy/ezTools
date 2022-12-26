package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/Anveena/ezTools/captcha/ezCaptchaPB"
	"github.com/Anveena/ezTools/networking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	ezCaptchaPB.UnimplementedEZCaptchaServiceServer
}

func startGRPCService() {
	listener, err := networking.ListenTCP(captchaConfig.GRPCPort, -1, -1)
	if err != nil {
		panic(err)
	}
	cert, err := tls.LoadX509KeyPair(captchaConfig.CrtPath, captchaConfig.KeyPath)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(grpc.Creds(credentials.NewServerTLSFromCert(&cert)))
	ezCaptchaPB.RegisterEZCaptchaServiceServer(s, &server{})
	if err = s.Serve(listener); err != nil {
		panic(fmt.Errorf("grpc服务启动失败:%s", err.Error()))
	}
}

func (s server) GetCaptcha(_ context.Context, _ *ezCaptchaPB.EZCaptchaEmpty) (*ezCaptchaPB.EZCaptchaRsp, error) {
	data, correctAnswer, err := gen()
	if err != nil {
		return &ezCaptchaPB.EZCaptchaRsp{
			Suc:     false,
			ErrDesc: err.Error(),
		}, nil
	}
	return &ezCaptchaPB.EZCaptchaRsp{
		Suc:           true,
		CorrectAnswer: correctAnswer,
		PngBase64:     base64.StdEncoding.EncodeToString(data),
	}, nil
}

package main

import (
	"fmt"
	"github.com/Anveena/ezTools/log/model"
	"github.com/Anveena/ezTools/networking"
	"google.golang.org/grpc"
	"runtime"
)

func startGRPCService() {
	lis, err := networking.ListenTCP(ezLSConfig.GRPCPort, -1, -1)
	if err != nil {
		panic(fmt.Errorf("监听tcp连接失败:%s", err.Error()))
	}
	s := grpc.NewServer()
	model.RegisterEzLogGrpcServer(s, &server{})
	if err = s.Serve(lis); err != nil {
		panic(fmt.Errorf("grpc服务启动失败:%s", err.Error()))
	}
}

type server struct {
	model.UnimplementedEzLogGrpcServer
}

func (s *server) Log(logServer model.EzLogGrpc_LogServer) error {
	runtime.LockOSThread()
	var err error
	for {
		lm := new(model.EZLogReq)
		err = logServer.RecvMsg(lm)
		if err != nil {
			println(err.Error())
			return nil
		}
		logModelChan <- lm
	}
}

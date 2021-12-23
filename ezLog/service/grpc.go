package main

import (
	"fmt"
	"github.com/Anveena/ezTools/ezLog/ezLogPB"
	"github.com/Anveena/ezTools/ezNetworking"
	"google.golang.org/grpc"
	"runtime"
)

func startGRPCService() {
	lis, err := ezNetworking.ListenTCP(ezLSConfig.GRPCPort, -1, -1)
	if err != nil {
		panic(fmt.Errorf("监听tcp连接失败:%s", err.Error()))
	}
	s := grpc.NewServer()
	ezLogPB.RegisterEzLogGrpcServer(s, &server{})
	if err = s.Serve(lis); err != nil {
		panic(fmt.Errorf("grpc服务启动失败:%s", err.Error()))
	}
}

type server struct {
	ezLogPB.UnimplementedEzLogGrpcServer
}

func (s *server) Log(logServer ezLogPB.EzLogGrpc_LogServer) error {
	runtime.LockOSThread()
	var err error
	for {
		lm := new(ezLogPB.EZLogReq)
		err = logServer.RecvMsg(lm)
		if err != nil {
			println(err.Error())
			return nil
		}
		logModelChan <- lm
	}
}

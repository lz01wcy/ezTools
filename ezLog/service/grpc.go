package main

import (
	"github.com/Anveena/ezTools/ezLog/ezLogPB"
	"github.com/Anveena/ezTools/ezNetworking"
	"google.golang.org/grpc"
	"log"
	"runtime"
	"time"
)

func startGRPCService() error {
	lis, err := ezNetworking.ListenTCP4(ezLSConfig.GRPCPort, -1, -1)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ezLogPB.RegisterEzLogGrpcServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}

type server struct {
	ezLogPB.UnimplementedEzLogGrpcServer
}

func (s *server) Log(logServer ezLogPB.EzLogGrpc_LogServer) error {
	runtime.LockOSThread()
	lm := new(ezLogPB.LogReq)
	var err error
	for {
		err = logServer.RecvMsg(lm)
		if err != nil {
			println(err.Error())
			return nil
		}
		logModelChan <- &logModel{
			Level:    lm.Level,
			AppName:  lm.AppName,
			FileName: lm.FileName,
			FileLine: lm.FileLine,
			Tag:      lm.Tag,
			Time:     time.UnixMicro(lm.Time),
			Content:  lm.Content,
		}
	}
}

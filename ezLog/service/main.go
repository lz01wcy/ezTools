package main

import (
	"github.com/Anveena/ezTools/ezConfig"
	"github.com/Anveena/ezTools/ezLog/ezLogPB"
	"github.com/Anveena/ezTools/ezMySQL"
)

type ezLogServiceConfig struct {
	HowManyDaysThatLogsShouldSave int
	HowManyLogsToInsertDBOnce     int
	HowOftenToInsertDBInSeconds   int
	LogModelChanSize              int
	GRPCPort                      int
	MySQLConf                     ezMySQL.Info
}

var ezLSConfig = &ezLogServiceConfig{}
var logModelChan chan *ezLogPB.EZLogReq

func main() {
	if err := ezConfig.ReadConf(ezLSConfig); err != nil {
		println(err.Error())
		return
	}
	logModelChan = make(chan *ezLogPB.EZLogReq, ezLSConfig.LogModelChanSize)
	go startDBWritingThread()
	startGRPCService()
}

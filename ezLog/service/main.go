package main

import (
	"github.com/Anveena/ezTools/ezConfig"
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
var logModelChan chan *logModel

func main() {
	if err := ezConfig.ReadConf(ezLSConfig); err != nil {
		println(err.Error())
		return
	}
	logModelChan = make(chan *logModel, ezLSConfig.LogModelChanSize)
	failed := make(chan bool, 10)
	go func() {
		defer func() {
			failed <- true
		}()
		if err := startDBWritingThread(); err != nil {
			println(err.Error())
		}
	}()
	go func() {
		defer func() {
			failed <- true
		}()
		if err := startGRPCService(); err != nil {
			println(err.Error())
		}
	}()
	<-failed
}

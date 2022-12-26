package main

import (
	"github.com/Anveena/ezTools/config"
	"github.com/Anveena/ezTools/log/model"
	"github.com/Anveena/ezTools/mysql"
)

type ezLogServiceConfig struct {
	HowManyDaysThatLogsShouldSave int
	HowManyLogsToInsertDBOnce     int
	HowOftenToInsertDBInSeconds   int
	LogModelChanSize              int
	GRPCPort                      int
	MySQLConf                     mysql.Info
}

var ezLSConfig = &ezLogServiceConfig{}
var logModelChan chan *model.EZLogReq

func (c *ezLogServiceConfig) Check() {

}
func main() {
	config.ReadConf(ezLSConfig)
	logModelChan = make(chan *model.EZLogReq, ezLSConfig.LogModelChanSize)
	go startDBWritingThread()
	startGRPCService()
}

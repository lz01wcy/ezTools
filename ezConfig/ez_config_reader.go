package ezConfig

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func ReadConf(configModel Checkable) {
	confPath := ""
	flag.StringVar(&confPath, "c", "", "Config toml File Path")
	flag.Parse()
	if confPath == "" {
		panic("配置文件路径没有输入.需要输入-c xx")
	}
	tomlData, err := os.ReadFile(confPath)
	if err != nil {
		panic(fmt.Sprintf("读取配置文件错误:%s", err.Error()))
	}
	if err = toml.Unmarshal(tomlData, configModel); err != nil {
		panic(fmt.Sprintf("toml解析错误:%s", err.Error()))
	}
	configModel.Check()
	return
}

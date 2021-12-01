package ezConfig

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

func ReadConf(configModel interface{}) error {
	confPath := ""
	flag.StringVar(&confPath, "c", "", "Config toml File Path")
	flag.Parse()
	if confPath == "" {
		return fmt.Errorf("配置文件路径没有输入.需要输入-c xx")
	}
	tomlData, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	return toml.Unmarshal(tomlData, configModel)
}

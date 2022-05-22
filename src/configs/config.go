package configs

import (
	"SqlLineage/src/configs/properties"
	"fmt"
	"github.com/spf13/viper"
)

type PropertyReader interface {
}

var path string = "settings/application.yml"

var readers map[string]PropertyReader = make(map[string]PropertyReader, 0)

func init() {
	ReadConfig()
	SetUpConfig()
}

func ReadConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("找不到配置文件 path:", path)
		return
	}
}

// 初始化所有配置文件
func SetUpConfig() {
	properties.SetMysql()
	properties.SetRedis()
	properties.SetServer()
}

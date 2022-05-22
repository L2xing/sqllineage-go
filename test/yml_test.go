package test

import (
	"github.com/spf13/viper"
	"testing"
)

func TestRead(t *testing.T) {
	path := "../settings/application.yml"
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		println("找不到配置文件 path:", path)
		return
	}
	get := viper.GetString("abc")
	println(get)
}

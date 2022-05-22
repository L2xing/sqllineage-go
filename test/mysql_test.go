package test

import (
	"github.com/spf13/viper"
	"testing"
)

func TestCon(t *testing.T) {
	path := "../settings/application.yml"
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		println("找不到配置文件 path:", path)
		return
	}

	//My := new(properties.MySQL)
	//sub := viper.Sub("mysql")
	//sub.Unmarshal(My)
	////通过访问结构体成员获取yaml文件中对应的key-value
	//fmt.Println(My.Host)
	//fmt.Println(My.Password)
}

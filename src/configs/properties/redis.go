package properties

import (
	"fmt"
	"github.com/spf13/viper"
)

// 前缀
var subRedis string = "redis"

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}

var GolbalRedis *Redis

func SetRedis() {
	GolbalRedis = new(Redis)
	v := viper.Sub(subRedis)
	v.Unmarshal(GolbalRedis)
	fmt.Println(subRedis + "配置读入完成")
}

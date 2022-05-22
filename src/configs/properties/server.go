package properties

import (
	"fmt"
	"github.com/spf13/viper"
)

// 前缀
var subServer string = "server"

type Server struct {
	Port string `yaml:"port"`
}

var GolbalServer *Server

func SetServer() {
	GolbalServer = new(Server)
	v := viper.Sub(subServer)
	v.Unmarshal(GolbalServer)
	fmt.Println(subServer + "配置读入完成")
}

package properties

import (
	"fmt"
	"github.com/spf13/viper"
)

// 前缀
var subMysql string = "mysql"

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}

var GolbalMySQL *Mysql

func SetMysql() {
	GolbalMySQL = new(Mysql)
	v := viper.Sub(subMysql)
	v.Unmarshal(GolbalMySQL)
	fmt.Println(subMysql + "配置读入完成")
}

package main

import (
	"SqlLineage/src/configs"
	"SqlLineage/src/controllers"
)

func main() {
	println("程序的入口")
	// 连接数据库
	configs.SetUpGrom()
	// Gin启动
	controllers.Run()
}

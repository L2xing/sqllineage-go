package controllers

import (
	"SqlLineage/src/models/bo"
	"SqlLineage/src/services/app"
	jwt "SqlLineage/src/utils/jwt"
	"SqlLineage/src/utils/net"
	"github.com/gin-gonic/gin"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, authRouter)
}

func authRouter(group *gin.RouterGroup) {
	auth := group.Group("/auth")
	auth.POST("/register", Register)
	auth.POST("/login", Login)
}

func Login(c *gin.Context) {
	// 用户发送用户名和密码过来
	var user bo.UserInfo
	err := c.BindJSON(&user)
	if err != nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	// 校验用户名和密码是否正确
	verify := app.Verify(user.Username, user.Password)
	if verify == nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	if token, err := jwt.GenToken(verify); err != nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	} else {
		net.ResponseSuccess(c, gin.H{
			"token": token,
		})
	}
}

func Register(c *gin.Context) {
	// 用户发送用户名和密码过来
	var user bo.UserInfo
	err := c.BindJSON(&user)
	if err != nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	verify := app.Register(user.Username, user.Password)
	if verify == nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	if token, err := jwt.GenToken(verify); err != nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	} else {
		net.ResponseSuccess(c, gin.H{
			"token": token,
		})
	}
}

//// token 获取
//func authHandler(c *gin.Context) {
//
//}
//
//func getUser(id int64) {
//
//}
//
//func addUser(user do.User) {
//
//}
//
//func deleteUser(id int64) {
//	c
//
//}

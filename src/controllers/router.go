package controllers

import (
	"SqlLineage/src/configs/properties"
	"SqlLineage/src/utils/jwt"
	"SqlLineage/src/utils/net"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const (
	currentUser     = "current-user"
	currentUserName = "current-user-username"
)

var (
	routerNoCheckRole = make([]func(*gin.RouterGroup), 0)
	routerCheckRole   = make([]func(v1 *gin.RouterGroup), 0)
)

// Setup 路由设置
func Run() {
	r := gin.Default()
	// 添加自定义异常处理
	r.Use(ErrorHandler)
	r.GET("/ping", func(context *gin.Context) {
		println("123")
	})
	// 添加访问路由映射
	InitAdminRouter(r)
	port := fmt.Sprintf(":%s", properties.GolbalServer.Port)
	r.Run(port)
}

//InitAdminRouter 后台模块路由
func InitAdminRouter(r *gin.Engine) *gin.Engine {
	// 无需认证的路由
	adminNoCheckRoleRouter(r)
	// 需要认证的路由
	adminCheckRoleRouter(r)
	return r
}

func adminNoCheckRoleRouter(r *gin.Engine) {
	// 可根据业务需求来设置接口版本
	v1 := r.Group("/api")
	// 空接口防止v1定义无使用报错
	v1.GET("/nilcheckrole", nil)

	for _, f := range routerNoCheckRole {
		f(v1)
	}
}

func adminCheckRoleRouter(r *gin.Engine) {
	// 可根据业务需求来设置接口版本
	v1 := r.Group("/api", JWTAuthMiddleware)
	// 空接口防止v1定义无使用报错
	v1.GET("/checkrole", nil)
	for _, f := range routerCheckRole {
		f(v1)
	}
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware(c *gin.Context) {
	//客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	//这里假设Token放在Header的Authorization中，并使用Bearer开头
	//这里的具体实现方式要依据你的实际业务情况决定
	authHeader := c.Request.Header.Get("Auth")
	if authHeader == "" {
		net.ResponseError(c, net.CodeLoginExpire)
		c.Abort()
		return
	}

	// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
	token, err := jwt.ParseToken(authHeader)
	if err != nil {
		logrus.Errorln("token解析失败")
		net.ResponseError(c, net.CodeLoginExpire)
		c.Abort()
		return
	}
	// 续签 依赖Redis
	c.Set(currentUser, token)
	c.Set(currentUserName, token.Username)

	c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
}

// 获取当前用户名
func GetOnlineUserName(c *gin.Context) (string, bool) {
	if get, exists := c.Get(currentUserName); exists {
		return get.(string), true
	} else {
		return "", false
	}
}

// 异常处理中间件
func ErrorHandler(c *gin.Context) {
	defer func() {
		err := recover()
		if err != nil {
			switch err.(type) {
			case net.ResponseData:
				data := err.(net.ResponseData)
				net.ResponseErrorWithData(c, &data)
			default:
				net.ResponseError(c, net.CodeSeverError)
			}
		}
	}()
	c.Next()
}

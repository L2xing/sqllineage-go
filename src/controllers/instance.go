package controllers

import (
	"SqlLineage/src/models/do"
	"SqlLineage/src/services/app"
	"SqlLineage/src/utils/net"
	"github.com/gin-gonic/gin"
	"strconv"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, instanceRouter)
	routerCheckRole = append(routerCheckRole, instanceAuthRouter)
}

func instanceRouter(group *gin.RouterGroup) {

}

func instanceAuthRouter(group *gin.RouterGroup) {
	instance := group.Group("/instance")
	instance.POST("", InsertInstance)
	instance.GET("", QueryInstances)
	instance.GET("/all/:id", PullInstanceData)
	instance.DELETE("/:id", DeleteInstance)
}

func InsertInstance(c *gin.Context) {
	var instance do.Instance
	err := c.BindJSON(&instance)
	if err != nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	if name, b := GetOnlineUserName(c); b {
		instance.Owner = name
	} else {
		net.ResponseError(c, net.CodeLoginExpire)
		return
	}
	if !app.AddInstance(&instance) {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	net.ResponseSuccess(c, nil)
}

func QueryInstances(c *gin.Context) {
	var page net.PageInfo
	var username string
	if info, err := net.GetPageInfo(c); err != nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	} else {
		page = info
	}
	if name, b := GetOnlineUserName(c); b {
		username = name
	} else {
		net.ResponseError(c, net.CodeLoginExpire)
		return
	}
	pageResult := app.QueryInstancesPage(page, username)
	net.ResponseSuccess(c, pageResult)
}

func PullInstanceData(c *gin.Context) {
	idStr := c.Param("id")
	atoi, err := strconv.Atoi(idStr)
	if err != nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	data := app.PullInstanceData(uint(atoi))
	net.ResponseSuccess(c, data)
}

func DeleteInstance(c *gin.Context) {
	idStr := c.Param("id")
	atoi, err := strconv.Atoi(idStr)
	if err != nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	name, _ := GetOnlineUserName(c)
	app.DeleteInstance(uint(atoi), name)
	net.ResponseSuccess(c, nil)
}

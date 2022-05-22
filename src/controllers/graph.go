package controllers

import (
	"SqlLineage/src/models/bo"
	"SqlLineage/src/services/app"
	"SqlLineage/src/utils/net"
	"github.com/gin-gonic/gin"
	"strconv"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, graphRouter)
	routerCheckRole = append(routerCheckRole, graphAuthRouter)
}

func graphRouter(group *gin.RouterGroup) {
	//routerGroup := group.Group("/user")
	//routerGroup.GET("/:id")
}

func graphAuthRouter(group *gin.RouterGroup) {
	graph := group.Group("/graph")
	graph.GET("/instances", GetInstancesGraph)
	graph.GET("/databases/:id", GetDatabasesGraph)
	graph.GET("/datasets/:id", GetDatasetsGraph)
	graph.GET("/columns/:id", GetColumnsGraph)
	graph.GET("/category", GetCategory)
}

func GetInstancesGraph(c *gin.Context) {
	var username string
	if name, b := GetOnlineUserName(c); b {
		username = name
	} else {
		net.ResponseError(c, net.CodeLoginExpire)
		return
	}
	nodes := app.GetInstanceByUsername(username)
	net.ResponseSuccess(c, nodes)
	return
}

func GetDatabasesGraph(c *gin.Context) {
	idStr := c.Param("id")
	atoi, err := strconv.Atoi(idStr)
	if err != nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	nodes := app.GetDatabaseByInstanceId(uint(atoi))
	net.ResponseSuccess(c, nodes)
	return
}

func GetDatasetsGraph(c *gin.Context) {
	idStr := c.Param("id")
	atoi, err := strconv.Atoi(idStr)
	if err != nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	nodes := app.GetDatasetsByDatabaseId(uint(atoi))
	net.ResponseSuccess(c, nodes)
	return
}

func GetColumnsGraph(c *gin.Context) {
	idStr := c.Param("id")
	atoi, err := strconv.Atoi(idStr)
	if err != nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	nodes := app.GetColumnsByDatasetId(uint(atoi))
	net.ResponseSuccess(c, nodes)
	return
}

func GetCategory(c *gin.Context) {
	net.ResponseSuccess(c, bo.GraphCategory)
	return
}

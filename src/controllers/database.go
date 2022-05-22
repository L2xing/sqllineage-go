package controllers

import (
	"SqlLineage/src/services/app"
	"SqlLineage/src/utils/net"
	"github.com/gin-gonic/gin"
	"strconv"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, databaseRouter)
	routerCheckRole = append(routerCheckRole, databaseAuthRouter)
}

func databaseRouter(group *gin.RouterGroup) {
	routerGroup := group.Group("/database")
	{
		routerGroup.PUT("")
	}
}

func databaseAuthRouter(group *gin.RouterGroup) {
	databaseGroup := group.Group("/database")
	databaseGroup.GET("/instance/:id", GetDatabases)
}

func GetDatabases(c *gin.Context) {
	instanceIdStr := c.Param("id")
	atoi, err := strconv.Atoi(instanceIdStr)
	if err != nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	instanceId := uint(atoi)
	databases := app.QueryDatabaseByInstanceId(instanceId)
	net.ResponseSuccess(c, databases)
}

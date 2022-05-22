package controllers

import (
	"SqlLineage/src/models/bo"
	"SqlLineage/src/services/app"
	"SqlLineage/src/services/sql"
	"SqlLineage/src/utils/net"
	"github.com/gin-gonic/gin"
	"sync"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, sqlRouter)
	routerCheckRole = append(routerCheckRole, sqlAuthRouter)
}

func sqlRouter(group *gin.RouterGroup) {
	//routerGroup := group.Group("/user")
	//routerGroup.GET("/:id")
}

func sqlAuthRouter(group *gin.RouterGroup) {
	sqlGroup := group.Group("/sql")
	sqlGroup.POST("/parse", PostSqlParse)
	sqlGroup.POST("/lineage", AddSqlLineage)
	sqlGroup.GET("/nodes", GetLineageNode)
}

var UserLineageMaps map[string]*sql.LineageContext = make(map[string]*sql.LineageContext, 0)
var mutex sync.RWMutex

func PostSqlParse(c *gin.Context) {
	var str bo.Sql
	c.BindJSON(&str)
	var sqlStr = str.Sql
	if sqlStr == "" {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	lineage, uuid := sql.ParserSqlLineage(sqlStr)
	if lineage == nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	cols := sql.BuildColsGraph(lineage)
	tables := sql.BuildTablesGraph(lineage)
	net.ResponseSuccess(c, gin.H{
		"tables": tables,
		"cols":   cols,
		"uuid":   uuid,
	})
}

func AddSqlLineage(c *gin.Context) {
	var sqlbo bo.AddSqlBo
	c.BindJSON(&sqlbo)
	if sqlbo.Uuid == "" {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	sqlLineage := sql.AddSqlLineage(&sqlbo)
	if sqlLineage {

	} else {

	}

	return
}

func GetGraphNode(c *gin.Context) {
	var query bo.SqlNodeQuery
	c.BindQuery(&query)
	if query.NodeId == 0 {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	node := app.GetGraphNode(query)
	net.ResponseSuccess(c, node)
}

func GetLineageNode(c *gin.Context) {
	var query bo.SqlNodeQuery
	c.BindQuery(&query)
	if query.NodeId == 0 {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	node := app.GetLineageNode(query.NodeId)
	net.ResponseSuccess(c, node)
}

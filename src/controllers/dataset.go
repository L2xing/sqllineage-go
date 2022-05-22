package controllers

import (
	"SqlLineage/src/services/app"
	"SqlLineage/src/utils/net"
	"github.com/gin-gonic/gin"
	"strconv"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, datasetRouter)
	routerCheckRole = append(routerCheckRole, datasetAuthRouter)
}

func datasetRouter(group *gin.RouterGroup) {

}

func datasetAuthRouter(group *gin.RouterGroup) {
	dataset := group.Group("/dataset")
	dataset.GET("/database/:id", GetDatasetsDatabaseId)
	dataset.GET("/vdatasets", GetVDatasets)
}

func GetDatasetsDatabaseId(c *gin.Context) {
	idStr := c.Param("id")
	atoi, err := strconv.Atoi(idStr)
	if err != nil {
		net.ResponseError(c, net.CodeBadRequest)
		return
	}
	datasets := app.QueryDatasetsByDatabaseId(uint(atoi))
	net.ResponseSuccess(c, datasets)
}

func GetVDatasets(c *gin.Context) {
	vDatasets := app.QueryVDatasets()
	net.ResponseSuccess(c, vDatasets)
}

package controllers

import "github.com/gin-gonic/gin"

func init() {
	routerNoCheckRole = append(routerNoCheckRole, columnRouter)
	routerCheckRole = append(routerCheckRole, columnAuthRouter)
}

func columnRouter(group *gin.RouterGroup) {

}

func columnAuthRouter(group *gin.RouterGroup) {

}

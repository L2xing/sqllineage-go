package controllers

import (
	"github.com/gin-gonic/gin"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, userRouter)
	routerCheckRole = append(routerCheckRole, userAuthRouter)
}

func userRouter(group *gin.RouterGroup) {
	//routerGroup := group.Group("/user")
	//routerGroup.GET("/:id")
}

func userAuthRouter(group *gin.RouterGroup) {

}

//func getUser(id int64) {
//
//}
//
//func addUser(user do.User) {
//
//}
//
//func deleteUser(id int64) {
//
//}

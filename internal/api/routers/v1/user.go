package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(router *gin.RouterGroup) {
	router.GET("/get", controllers.User.GetUser)
	router.POST("/register", controllers.User.Register)
	router.PUT("/update", controllers.User.UpdateUsers)
	router.POST("/login", controllers.User.Login)
}

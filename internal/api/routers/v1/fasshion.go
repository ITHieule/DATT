package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterFashionRouter(router *gin.RouterGroup) {
	router.GET("/get", controllers.Fashion.Get)
}

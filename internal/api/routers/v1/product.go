package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterProductRouter(router *gin.RouterGroup) {
	router.GET("/get", controllers.Product.Getproduct)
	router.GET("/getimage", controllers.Product.Getproduct_image)
}

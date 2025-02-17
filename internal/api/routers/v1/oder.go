package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterOderRouter(router *gin.RouterGroup) {
	router.GET("/get", controllers.Oder.Getoder)
	router.GET("/getorderid", controllers.Oder.GetProductDetailsByOrderID)
	
}
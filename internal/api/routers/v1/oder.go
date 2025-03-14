package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterOderRouter(router *gin.RouterGroup) {
	router.POST("/get", controllers.Oder.GetOder)
	router.GET("/getorderid/:order_id", controllers.Oder.GetProductDetailsByOrderID)
	router.POST("/create", controllers.Oder.CreateOrderFromCart)
}

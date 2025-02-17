package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCartRouter(router *gin.RouterGroup) {
	router.GET("/get", controllers.Cart.GetToCart)
	router.POST("/addtocart", controllers.Cart.AddToCart)
	router.PUT("/update", controllers.Cart.UpdateCartQuantity)
	router.DELETE("/delete", controllers.Cart.RemoveFromCart)
}

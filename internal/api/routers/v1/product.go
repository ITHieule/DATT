package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterProductRouter(router *gin.RouterGroup) {
	router.GET("/get", controllers.Product.Getproduct)
	router.POST("/add", controllers.Product.AddProductController)
	router.PUT("/update", controllers.Product.UpdateProductController)
	router.DELETE("/delete", controllers.Product.Deleteproduct)
	router.POST("/search", controllers.Product.SearchProduct)

	router.GET("/getimage", controllers.Product.Getproduct_image)
	router.GET("/getvariants/:id", controllers.Product.GetProductDetailController)

}

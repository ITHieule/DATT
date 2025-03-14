package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAddressRouter(router *gin.RouterGroup) {
	router.POST("/get", controllers.Address.GetAddressByUserID)
	router.POST("/create", controllers.Address.CreateAddressByUserID)
	router.PUT("/update", controllers.Address.UpdateAddressByUserID)
	router.DELETE("/delete", controllers.Address.DeleteAddressByUserID)
}

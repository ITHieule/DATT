package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCaterogyRouter(router *gin.RouterGroup) {
	router.GET("/get", controllers.Caterogy.GetCaterogy)
}

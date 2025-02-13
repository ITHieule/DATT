package router_v1

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	RegisterCommonRouter(v1.Group(""))

	RegisterReportRouter(v1.Group("/rp"))

	RegisterAdminRouter(v1.Group("/rp/admin"))

	RegisterFashionRouter(v1.Group("/fashion"))

	RegisterUserRouter(v1.Group("/user"))

	RegisterCartRouter(v1.Group("/cart"))
}

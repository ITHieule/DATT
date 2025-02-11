package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

// FILE - ADMIN ROUTER
func RegisterAdminRouter(router *gin.RouterGroup) {

	//api - Login
	router.POST("/login", controllers.Admin.Login)

	//api - GetReportByIsAnsweredAndRpid
	router.POST("/isanswered", controllers.Admin.GetReportByIsAnsweredAndRpid)

	//api - SendAnswer
	router.POST("/sendanswer", controllers.Admin.SendAnswer)

	//api - UpdateAnswer
	router.PATCH("/updateanswer", controllers.Admin.UpdateAnswer)

	//api - UpdateRPType
	router.PATCH("/updaterptype", controllers.Admin.UpdateRPTypes)

	//api - GetRPTypes
	router.GET("/rptypes", controllers.Admin.GetRPTypes)

	//api - ExcelByDepartments
	router.GET("/excelbydepartments", controllers.Admin.ExcelByDepartments)

	//api - ExcelByReportTypes
	router.GET("/excelbyreporttypes", controllers.Admin.ExcelByReportTypes)

}

package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

// FILE - REPORT ROUTER
func RegisterReportRouter(router *gin.RouterGroup) {
	//api - Get All Reports
	router.GET("/reports", controllers.Rp.GetReport)

	//api - Get Report By Report Code
	router.POST("/rpbycode", controllers.Rp.GetReportByRpcode)

	//api - Insert Report
	router.POST("/insertreport1", controllers.Rp.InsertFirstReport)

	//api - Check And Get Report
	router.POST("/checkfeedback", controllers.Rp.CheckRPCodeAndGetReports)

	//api - Send Report
	router.POST("/sendfeedback", controllers.Rp.SendReport)

	//api - Update Report
	router.PATCH("/updatenote", controllers.Rp.UpdateNote)

	//api - Update Is Completed Report
	router.PATCH("/iscompleted", controllers.Rp.UpdateIsCompleted)

	//api - Get Work Units
	router.GET("/departments", controllers.Rp.GetDepartments)

	//api - Send Email
	router.POST("/sendemail", controllers.Rp.SendEmail)

	

}

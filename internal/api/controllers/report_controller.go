package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

// FILE - REPORT CONTROLLER
type ReportController struct {
	*BaseController
}

var Rp = &ReportController{}

// func GetReportController
func (c *ReportController) GetReport(ctx *gin.Context) {

	// gọi func GetReportService từ service ko cần tham số
	result, err := services.RP.GetReportService()

	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// Trả về kết quả json
	ctx.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

// func GetReportByRpcodeController
func (c *ReportController) GetReportByRpcode(ctx *gin.Context) {
	var requestParams request.GetReportByRpcodeRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.RP.GetReportByRpcodeService(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

// func InsertReportController
func (c *ReportController) InsertFirstReport(ctx *gin.Context) {
	var requestParams request.InsertReportRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.RP.InsertReportFirstService(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// func - CheckRPCodeAndGetReportsController
func (c *ReportController) CheckRPCodeAndGetReports(ctx *gin.Context) {
	var requestParams request.GetReportByRpcodeRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.RP.CheckRPCodeAndGetReportsService(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// func - SendReportController
func (c *ReportController) SendReport(ctx *gin.Context) {
	var requestParams request.InsertReportRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.RP.SendReportService(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// func - GetWorkUnitController
func (c *ReportController) GetDepartments(ctx *gin.Context) {
	result, err := services.RP.GetDepartmentsService()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// func - SendEmailController
func (c *ReportController) SendEmail(ctx *gin.Context) {
	// Gọi func SendEmailService cơ sở dữ liệu
	result, err := services.RP.SendEmailService()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// func - UpdateNoteController
func (c *ReportController) UpdateNote(ctx *gin.Context) {
	var requestParams request.UpdateNoteRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.RP.UpdateNoteService(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// func - UpdateIsCompletedController
func (c *ReportController) UpdateIsCompleted(ctx *gin.Context) {
	var requestParams request.GetReportByRpcodeRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.RP.UpdateIsCompletedService(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

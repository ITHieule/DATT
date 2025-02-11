package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

// FILE - ADMIN CONTROLLER
type AdminController struct {
	*BaseController
}

var Admin = &AdminController{}

// func LoginController
func (c *AdminController) Login(ctx *gin.Context) {
	var requestParams request.LoginRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	token, err := services.Admin.LoginService(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, err.Error())
		return
	}

	response.OkWithData(ctx, gin.H{"token": token})
}

// func - GetReportByIsAnsweredAndRpidService
func (c *AdminController) GetReportByIsAnsweredAndRpid(ctx *gin.Context) {
	var requestParams request.GetReportByIsAnsweredAndRpidRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.Admin.GetReportByIsAnsweredAndRpidService(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// func - SendAnswerController
func (c *AdminController) SendAnswer(ctx *gin.Context) {
	var requestParams request.SendAnswerRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.Admin.SendAnswerService(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// func - UpdateAnswerController
func (c *AdminController) UpdateAnswer(ctx *gin.Context) {
	var requestParams request.UpdateAnswerRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.Admin.UpdateAnswerService(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// func - GetRPTypesController
func (c *AdminController) GetRPTypes(ctx *gin.Context) {
	result, err := services.Admin.GetRPTypesService()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// func - UpdateRPTypesController
func (c *AdminController) UpdateRPTypes(ctx *gin.Context) {
	var requestParams request.UpdateRPTypeRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.Admin.UpdateRPTypeService(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// func - ExcelByDepartmentsController
func (c *AdminController) ExcelByDepartments(ctx *gin.Context) {
	result, err := services.Admin.ExcelByDepartmentsService()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// func - ExcelByReportTypesController
func (c *AdminController) ExcelByReportTypes(ctx *gin.Context) {
	result, err := services.Admin.ExcelByReportTypesService()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	*BaseController
}

var User = &UserController{}

func (c *UserController) Register(ctx *gin.Context) {
	var requestParams request.CreateUserRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.User.Register(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
	}
	response.OkWithData(ctx, result)
}

func (c *UserController) UpdateUsers(ctx *gin.Context) {
	var requestParams request.CreateUserRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.User.UpdateUserSevice(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

func (c *UserController) Login(ctx *gin.Context) {
	var requestParams request.LoginRequests
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	token, err := services.User.Login(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, err.Error()) // Trả về 401 nếu đăng nhập thất bại
		return
	}

	response.OkWithData(ctx, token)
}
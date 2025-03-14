package controllers

import (
	"log"
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/response"
	"web-api/internal/pkg/models/types"

	"github.com/gin-gonic/gin"
)

type AddressController struct{}

var Address = &AddressController{}

func (c *AddressController) GetAddressByUserID(ctx *gin.Context) {
	var requestParams types.ShippingAddress
	if err := ctx.ShouldBindJSON(&requestParams); err != nil {
		log.Printf("Invalid JSON data: %v", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid JSON data")
		return
	}

	result, err := services.Address.GetAddressByUserID(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *AddressController) CreateAddressByUserID(ctx *gin.Context) {
	var requestParams types.ShippingAddress

	// Bind JSON từ body request
	if err := ctx.ShouldBindJSON(&requestParams); err != nil {
		log.Printf("Invalid JSON data: %v", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid JSON data")
		return
	}

	// Gọi service để tạo địa chỉ mới
	result, err := services.Address.CreateAddressByUserID(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result) // Trả về địa chỉ vừa tạo
}

func (c *AddressController) UpdateAddressByUserID(ctx *gin.Context) {
	var requestParams types.ShippingAddress
	if err := ctx.ShouldBindJSON(&requestParams); err != nil {
		log.Printf("Invalid JSON data: %v", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid JSON data")
		return
	}

	err := services.Address.UpdateAddressByUserID(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithMessage(ctx, "Address updated successfully")
}

func (c *AddressController) DeleteAddressByUserID(ctx *gin.Context) {
	var requestParams types.ShippingAddress
	if err := ctx.ShouldBindJSON(&requestParams); err != nil {
		log.Printf("Invalid JSON data: %v", err)
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid JSON data")
		return
	}

	if err := services.Address.DeleteAddressByUserID(&requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

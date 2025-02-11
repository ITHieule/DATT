package controllers

import (
	"fmt"
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type FashionController struct {
	*BaseController
}

var Fashion = &FashionController{}

func (c *FashionController) Get(ctx *gin.Context) {
	result, err := services.FashionBusiness.FashionSevice()
	if err != nil {
		// Log lỗi chi tiết để dễ dàng kiểm tra trong hệ thống
		fmt.Println("Error in GetEnergySevice:", err)
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, "Failed to get energy records: "+err.Error())
		return
	}

	// Kiểm tra nếu kết quả trả về rỗng
	if len(result) == 0 {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "No energy records found.")
		return
	}

	response.OkWithData(ctx, result)
}

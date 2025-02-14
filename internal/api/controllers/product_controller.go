package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	*BaseController
}

var Product = &ProductController{}

func (c *ProductController) Getproduct(ctx *gin.Context) {
	result, err := services.ProductService.ProductSevice()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}
func (c *ProductController) Getproduct_image(ctx *gin.Context) {
	result, err := services.ProductService.Product_imageSevice()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// GetProductDetailController lấy thông tin chi tiết sản phẩm theo ID
func (c *ProductController) GetProductDetailController(ctx *gin.Context) {
	// Parse dữ liệu từ body
	var req request.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid request body")
		return
	}

	// Gọi service để lấy thông tin chi tiết sản phẩm
	product, err := services.ProductService.GetProductByID(req.Id)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, err.Error())
		return
	}

	// Trả về kết quả
	response.OkWithData(ctx, product)
}

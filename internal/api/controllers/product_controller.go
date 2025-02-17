package controllers

import (
	"net/http"
	"strconv"
	"web-api/internal/api/services"
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

// GetProductDetailController lấy thông tin chi tiết sản phẩm theo ID từ URL
func (c *ProductController) GetProductDetailController(ctx *gin.Context) {
	// Lấy id từ URL params
	productId := ctx.Param("id") // Tham số 'id' trong URL

	// Chuyển id từ chuỗi sang kiểu số nếu cần (ví dụ: nếu id là số nguyên)
	// Nếu id không hợp lệ, trả về lỗi
	id, err := strconv.Atoi(productId)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid product ID")
		return
	}

	// Gọi service để lấy thông tin chi tiết sản phẩm
	product, err := services.ProductService.GetProductByID(id)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, err.Error())
		return
	}

	// Trả về kết quả
	response.OkWithData(ctx, product)
}

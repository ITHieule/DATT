package controllers

import (
	"net/http"
	"strconv"
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

// AddProductController thêm sản phẩm cùng với biến thể
func (c *ProductController) AddProductController(ctx *gin.Context) {
	var req request.CreateProductRequest

	// Parse dữ liệu từ body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid request body")
		return
	}

	// Gọi service để thêm sản phẩm
	product, err := services.ProductService.AddProductService(&req)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// Trả về kết quả
	response.OkWithData(ctx, product)
}

func (c *ProductController) UpdateProductController(ctx *gin.Context) {
	var req request.CreateProductRequest

	// Parse dữ liệu từ body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid request body")
		return
	}

	// Gọi service để cập nhật sản phẩm
	product, err := services.ProductService.UpdateProductService(&req)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// Trả về kết quả
	response.OkWithData(ctx, product)
}

func (c *ProductController) Deleteproduct(ctx *gin.Context) {
	var requestParams request.CreateProductRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	err := services.ProductService.DeleteproductSevice(requestParams.Id)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, nil)
}

func (c *ProductController) SearchProduct(ctx *gin.Context) {
	var requestParams request.CreateProductRequest

	// Lấy tham số từ query string hoặc body
	if err := ctx.ShouldBindJSON(&requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid request body")
		return
	}

	// Gọi service để tìm kiếm sản phẩm
	result, err := services.ProductService.SearchProductService(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// Trả về dữ liệu JSON
	response.OkWithData(ctx, result)
}

func (c *ProductController) GetLatestProductHots(ctx *gin.Context) {
	// Lấy số lượng sản phẩm từ query params (mặc định lấy 12 sản phẩm)
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid limit value")
		return
	}

	// Gọi service để lấy danh sách sản phẩm hot mới nhất
	products, err := services.ProductService.GetLatestProductHots(limit)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// Trả về kết quả
	response.OkWithData(ctx, products)
}

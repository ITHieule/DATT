package controllers

import (
	"net/http"

	"fmt"
	"github.com/gin-gonic/gin"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/response"
)

type CartController struct {
	*BaseController
}

var Cart = &CartController{}

func (c *CartController) GetToCart(ctx *gin.Context) {
	// Lấy user_id từ query parameter
	userID := ctx.DefaultQuery("user_id", "")

	// Kiểm tra user_id có hợp lệ không
	if userID == "" {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "User ID is required")
		return
	}

	// Chuyển đổi user_id từ string sang int
	var userIDInt int
	if _, err := fmt.Sscanf(userID, "%d", &userIDInt); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid User ID")
		return
	}

	// Gọi service để lấy giỏ hàng theo userID
	result, err := services.Cart.GetCartByUserID(userIDInt)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// Trả về kết quả thành công
	response.OkWithData(ctx, result)
}

func (c *CartController) AddToCart(ctx *gin.Context) {
	var req struct {
		UserID           int `json:"user_id" binding:"required"`
		ProductVariantID int `json:"product_variant_id" binding:"required"`
		Quantity         int `json:"quantity" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.Cart.AddToCart(req.UserID, req.ProductVariantID, req.Quantity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully"})
}

func (c *CartController) UpdateCartQuantity(ctx *gin.Context) {
	var req struct {
		UserID           int `json:"user_id" binding:"required"`
		ProductVariantID int `json:"product_variant_id" binding:"required"`
		Quantity         int `json:"quantity" binding:"required,min=1"`
	}

	// Kiểm tra request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gọi service để cập nhật số lượng
	err := services.Cart.UpdateCartQuantity(req.UserID, req.ProductVariantID, req.Quantity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cập nhật số lượng sản phẩm thành công"})
}

func (c *CartController) RemoveFromCart(ctx *gin.Context) {
	var req struct {
		UserID           int `json:"user_id" binding:"required"`
		ProductVariantID int `json:"product_variant_id" binding:"required"`
	}

	// Kiểm tra request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gọi service để xóa sản phẩm khỏi giỏ hàng
	err := services.Cart.RemoveFromCart(req.UserID, req.ProductVariantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Xóa sản phẩm khỏi giỏ hàng thành công"})
}

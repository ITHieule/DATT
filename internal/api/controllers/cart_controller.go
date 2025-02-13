package controllers

import (
	"net/http"
	"web-api/internal/api/services"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	*BaseController
}

var Cart = &CartController{}

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

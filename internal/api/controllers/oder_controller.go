package controllers

import (
	"net/http"
	"strconv"
	"web-api/internal/api/services"

	"web-api/internal/pkg/models/response"
	"web-api/internal/pkg/models/types"

	"github.com/gin-gonic/gin"
)

type OderController struct {
	*BaseController
}

var Oder = &OderController{}

func (c *OderController) GetOder(ctx *gin.Context) {
	var request struct {
		UserID int `json:"user_id"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid request body")
		return
	}

	if request.UserID == 0 {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "User ID is required")
		return
	}

	result, err := services.Order.GetOderByUserID(request.UserID)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}
func (c *OderController) GetProductDetailsByOrderID(ctx *gin.Context) {
	orderID := ctx.Param("order_id")

	parsedOrderID, err := strconv.Atoi(orderID)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid Order ID")
		return
	}

	result, err := services.Order.GetOrderByID(parsedOrderID)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, "Failed to get product details: "+err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// 📌 **API: Đặt hàng từ giỏ hàng của user**
func (c *OderController) CreateOrderFromCart(ctx *gin.Context) {
	// 🏷️ **Lấy user_id từ body request** 
	var request struct {
		UserID         int    `json:"user_id"`
		RecipientName  string `json:"recipient_name"`
		RecipientPhone string `json:"recipient_phone"`
	}

	// 📌 **Kiểm tra dữ liệu đầu vào**
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid request body")
		return
	}

	if request.UserID == 0 {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "User ID is required")
		return
	}

	if request.RecipientName == "" || request.RecipientPhone == "" {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Recipient name and phone are required")
		return
	}

	// 🛍️ **Tạo đơn hàng mới**
	var newOrder types.Order
	newOrder.RecipientName = request.RecipientName
	newOrder.RecipientPhone = request.RecipientPhone

	err := services.Order.CreateOrderFromCart(request.UserID, &newOrder)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// ✅ **Trả về kết quả thành công**
	response.OkWithData(ctx, gin.H{
		"order_id":        newOrder.ID,
		"recipient_name":  newOrder.RecipientName,
		"recipient_phone": newOrder.RecipientPhone,
		"total_price":     newOrder.TotalPrice,
		"status":          newOrder.Status,
	})
}

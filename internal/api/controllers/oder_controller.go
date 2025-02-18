package controllers

import (
	"net/http"
	"strconv"
	"web-api/internal/api/services"

	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type OderController struct {
	*BaseController
}

var Oder = &OderController{}

func (c *OderController) GetOder(ctx *gin.Context) {
	// Struct để nhận dữ liệu từ body
	var request struct {
		UserID int `json:"user_id"`
	}

	// Bind JSON từ body vào struct
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid request body")
		return
	}

	// Kiểm tra userID có hợp lệ không
	if request.UserID == 0 {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "User ID is required")
		return
	}

	// Gọi service để lấy giỏ hàng theo userID
	result, err := services.OderServi.GetOderByUserID(request.UserID)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// Trả về kết quả thành công
	response.OkWithData(ctx, result)
}
func (c *OderController) GetProductDetailsByOrderID(ctx *gin.Context) {
	// Lấy order_id từ URL parameters
	orderID := ctx.Param("order_id")

	// Convert orderID nếu cần (ví dụ: từ string sang int)
	parsedOrderID, err := strconv.Atoi(orderID)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid Order ID")
		return
	}

	// Gọi service để lấy sản phẩm theo order_id
	result, err := services.OderServi.GetOrderByID(parsedOrderID)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, "Failed to get product details: "+err.Error())
		return
	}

	// Trả về kết quả thành công
	response.OkWithData(ctx, result)
}

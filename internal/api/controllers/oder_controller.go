package controllers

import (
	"net/http"

	"web-api/internal/api/services"

	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type OderController struct {
	*BaseController
}

var Oder = &OderController{}

//	func (c *OderController) Getoder(ctx *gin.Context) {
//		result, err := services.OderServi.GetoderSevice()
//		if err != nil {
//			response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
//			return
//		}
//		response.OkWithData(ctx, result)
//	}
func (c *OderController) Getoder(ctx *gin.Context) {
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
	result, err := services.OderServi.GetCartByUserID(request.UserID)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// Trả về kết quả thành công
	response.OkWithData(ctx, result)
}

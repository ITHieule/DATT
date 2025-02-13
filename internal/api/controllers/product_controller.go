package controllers

import (
	"net/http"
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
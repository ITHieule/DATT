package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type CaterogyController struct {
	*BaseController
}

var Caterogy = &CaterogyController{}

func (c *CaterogyController) GetCaterogy(ctx *gin.Context) {

		result, err := services.Caterogy.GetCaterogySevice()
		if err != nil {
			response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
			return
		}
		response.OkWithData(ctx, result)
	
	
}

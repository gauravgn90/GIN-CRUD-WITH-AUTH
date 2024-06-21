package controller

import (
	"gauravgn90/gin-crud-with-auth/v2/model"

	"github.com/gin-gonic/gin"
)

type ProductController interface {
	SaveProduct(ctx *gin.Context) (model.Product, int, error)
	UpdateProduct(ctx *gin.Context) error
	DeleteProduct(ctx *gin.Context) error
	FindAll(ctx *gin.Context) ([]model.Product, error)
}

type ProductControllerImpl struct {
	
}

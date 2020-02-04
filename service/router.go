package service

import "github.com/gin-gonic/gin"

type RouterInterface interface {
	ApplyRoutes(rg *gin.RouterGroup)
}
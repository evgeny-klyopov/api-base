package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseApiInterface interface {
	Run(relativePath string, rg *gin.RouterGroup, controllers []Controller)
	ping(c *gin.Context)
}

type BaseApi struct {}

func (b BaseApi) Run(relativePath string, rg *gin.RouterGroup, controllers []Controller) {
	route := rg.Group(relativePath)
	{
		route.GET("/ping", b.ping)

		for _, controller := range controllers {
			controller.ApplyRoutes(route)
		}
	}

}

func (b *BaseApi) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


type ServicesInterface interface {
	ResourceInterface

	ValidateError(c *gin.Context, err error, code int)
	ResponseError(c *gin.Context, err error)
	ValidateCode(c *gin.Context, code int)
}

type Service struct {
	Resource
}

func NewService() ServicesInterface {
	var service ServicesInterface

	service = &Service{}

	return service
}

func (s *Service) ValidateError(c *gin.Context, err error, code int) {
	ginError := c.Error(err)
	c.AbortWithStatusJSON(code, ginError.JSON())
}

func (s *Service) ValidateCode(c *gin.Context, code int) {
	c.AbortWithStatus(code)
}

func (s *Service) ResponseError(c *gin.Context, err error) {
	ginError := c.Error(err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, ginError.JSON())
}



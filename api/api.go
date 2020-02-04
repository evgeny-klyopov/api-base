package api

import (
	v1 "api-base/api/v1"
	"api-base/service"
	"github.com/gin-gonic/gin"
)

type ExtendApiInterface interface {
	service.RouterInterface
	service.BaseApiInterface
	apiV1(rg *gin.RouterGroup)
}

type Api struct {
	service.BaseApi
	Router  *gin.Engine
	Service service.ServicesInterface
}

func NewApi(r *gin.Engine) ExtendApiInterface {
	var api ExtendApiInterface

	api = &Api{
		BaseApi: service.BaseApi{},
		Service: service.NewService(),
		Router:  r,
	}

	return api
}

func (a *Api) ApplyRoutes(rg *gin.RouterGroup) {
	api := a.Router.Group("/api")
	{
		a.apiV1(api)
	}
}

func (a *Api) apiV1(rg *gin.RouterGroup) {
	a.Run(
		"v1",
		rg,
		[]service.Controller{
			v1.NewAuth(a.Service),
			v1.NewPosts(a.Service),
		},
	)
}

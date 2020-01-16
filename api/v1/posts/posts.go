package posts

import (
	"api-base/lib/middlewares"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Posts struct {
	service Service
}

type Service interface {
	Add(context.Context) error
	View(context.Context) error
	List(context.Context) error
	Remove(context.Context) error
	Update(context.Context) error
}


// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	p := &Posts{}
	posts := r.Group("/posts")
	{
		posts.GET("/", p.List)
		posts.POST("/add", middlewares.Authorized, p.Add)
		posts.GET("/:id", p.View)
		posts.DELETE("/:id", middlewares.Authorized, p.Remove)
		posts.PATCH("/:id", middlewares.Authorized, p.Update)


	}
}

func (p *Posts) validateError(c *gin.Context, err error, code int) {
	ginError := c.Error(err)
	c.AbortWithStatusJSON(code, ginError.JSON())
}

func (p *Posts) responseError(c *gin.Context, err error) {
	ginError := c.Error(err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, ginError.JSON())
}

package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)


type Auth struct {
	service Service
}

type Service interface {
	Login(context.Context) error
	Register(context.Context) error
	Check(context.Context) error
}


// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	a := &Auth{}
	auth := r.Group("/auth")
	{
		auth.POST("/login", a.Login)
		auth.POST("/sign-up", a.SignUp)
		auth.POST("/check", a.Check)
	}
}

func (a *Auth) validateError(c *gin.Context, err error, code int) {
	ginError := c.Error(err)
	c.AbortWithStatusJSON(code, ginError.JSON())
}

func (a *Auth) responseError(c *gin.Context, err error) {
	ginError := c.Error(err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, ginError.JSON())
}

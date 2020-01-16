package v1

import (
	"api-base/api/v1/auth"
	"api-base/api/v1/posts"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}


// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	{
		v1.GET("/ping", ping)
		auth.ApplyRoutes(v1)
		posts.ApplyRoutes(v1)
	}
}

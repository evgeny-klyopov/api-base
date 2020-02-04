package v1

import (
	"api-base/database/models"
	"api-base/lib/common"
	"api-base/lib/middlewares"
	"api-base/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Posts struct {
	service.ServicesInterface
}

func NewPosts(s service.ServicesInterface) service.Controller {
	var posts service.Controller

	posts = &Posts{s}

	return posts
}

func (p Posts) ApplyRoutes(rg *gin.RouterGroup) {
	posts := rg.Group("/posts")
	{
		posts.GET("/", p.list)
		posts.POST("/add", middlewares.Authorized, p.add)
		posts.GET("/:id", p.view)
		posts.DELETE("/:id", middlewares.Authorized, p.remove)
		posts.PATCH("/:id", middlewares.Authorized, p.update)

	}
}

func (p *Posts) add(c *gin.Context) {
	type RequestBody struct {
		Text string `json:"text" binding:"required"`
	}
	var requestBody RequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		p.ValidateError(c, err, http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(models.User)
	post := models.Post{Text: requestBody.Text, User: user}
	p.GetDB().NewRecord(post)
	p.GetDB().Create(&post)

	c.JSON(http.StatusOK, post.Serialize())

}
func (p *Posts) list(c *gin.Context) {
	cursor := c.Query("cursor")
	recent := c.Query("recent")

	var posts []models.Post

	if cursor == "" {
		if err := p.GetDB().Preload("User").Limit(10).Order("id desc").Find(&posts).Error; err != nil {
			p.ResponseError(c, err)
			return
		}
	} else {
		condition := "id < ?"
		if recent == "1" {
			condition = "id > ?"
		}
		if err := p.GetDB().Preload("User").Limit(10).Order("id desc").Where(condition, cursor).Find(&posts).Error; err != nil {
			p.ResponseError(c, err)
			return
		}
	}

	length := len(posts)
	serialized := make([]common.JSON, length, length)

	for i := 0; i < length; i++ {
		serialized[i] = posts[i].Serialize()
	}

	c.JSON(http.StatusOK, serialized)
}
func (p *Posts) view(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	// auto preloads the related model
	// http://gorm.io/docs/preload.html#Auto-Preloading
	if err := p.GetDB().Set("gorm:auto_preload", true).Where("id = ?", id).First(&post).Error; err != nil {
		p.ValidateError(c, err, http.StatusNotFound)

		return
	}

	c.JSON(http.StatusOK, post.Serialize())
}

func (p *Posts) remove(c *gin.Context) {
	id := c.Param("id")

	user := c.MustGet("user").(models.User)

	var post models.Post
	if err := p.GetDB().Where("id = ?", id).First(&post).Error; err != nil {
		p.ValidateError(c, err, http.StatusNotFound)
		return
	}

	if post.UserID != user.ID {
		p.ValidateError(c, errors.New("forbidden"), http.StatusForbidden)
		return
	}

	p.GetDB().Delete(&post)
	c.Status(http.StatusNoContent)
}

func (p *Posts) update(c *gin.Context) {
	id := c.Param("id")

	user := c.MustGet("user").(models.User)

	type RequestBody struct {
		Text string `json:"text" binding:"required"`
	}

	var requestBody RequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		p.ValidateError(c, err, http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := p.GetDB().Preload("User").Where("id = ?", id).First(&post).Error; err != nil {
		p.ValidateError(c, err, http.StatusNotFound)
		return
	}

	if post.UserID != user.ID {
		p.ValidateError(c, errors.New("forbidden"), http.StatusForbidden)
		return
	}

	post.Text = requestBody.Text
	p.GetDB().Save(&post)
	c.JSON(http.StatusOK, post.Serialize())
}
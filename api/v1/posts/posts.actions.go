package posts

import (
	"api-base/database/models"
	"api-base/lib/common"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func (p *Posts) Add(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	type RequestBody struct {
		Text string `json:"text" binding:"required"`
	}
	var requestBody RequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		p.validateError(c, err, http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(models.User)
	post := models.Post{Text: requestBody.Text, User: user}
	db.NewRecord(post)
	db.Create(&post)

	c.JSON(http.StatusOK, post.Serialize())

}
func (p *Posts) List(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	cursor := c.Query("cursor")
	recent := c.Query("recent")

	var posts []models.Post

	if cursor == "" {
		if err := db.Preload("User").Limit(10).Order("id desc").Find(&posts).Error; err != nil {
			p.responseError(c, err)
			return
		}
	} else {
		condition := "id < ?"
		if recent == "1" {
			condition = "id > ?"
		}
		if err := db.Preload("User").Limit(10).Order("id desc").Where(condition, cursor).Find(&posts).Error; err != nil {
			p.responseError(c, err)
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
func (p *Posts) View(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var post models.Post

	// auto preloads the related model
	// http://gorm.io/docs/preload.html#Auto-Preloading
	if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&post).Error; err != nil {
		p.validateError(c, err, http.StatusNotFound)

		return
	}

	c.JSON(http.StatusOK, post.Serialize())
}

func (p *Posts) Remove(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	user := c.MustGet("user").(models.User)

	var post models.Post
	if err := db.Where("id = ?", id).First(&post).Error; err != nil {
		p.validateError(c, err, http.StatusNotFound)
		return
	}

	if post.UserID != user.ID {
		p.validateError(c, errors.New("forbidden"), http.StatusForbidden)
		return
	}

	db.Delete(&post)
	c.Status(http.StatusNoContent)
}

func (p *Posts) Update(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	user := c.MustGet("user").(models.User)

	type RequestBody struct {
		Text string `json:"text" binding:"required"`
	}

	var requestBody RequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		p.validateError(c, err, http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := db.Preload("User").Where("id = ?", id).First(&post).Error; err != nil {
		p.validateError(c, err, http.StatusNotFound)
		return
	}

	if post.UserID != user.ID {
		p.validateError(c, errors.New("forbidden"), http.StatusForbidden)
		return
	}

	post.Text = requestBody.Text
	db.Save(&post)
	c.JSON(http.StatusOK, post.Serialize())
}
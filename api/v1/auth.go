package v1

import (
	"api-base/database/models"
	"api-base/lib/common"
	"api-base/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Auth struct {
	service.ServicesInterface
}

func NewAuth(s service.ServicesInterface) service.Controller {
	var auth service.Controller

	auth = &Auth{s}

	return auth
}

func (a Auth) ApplyRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", a.login)
		auth.POST("/sign-up", a.signUp)
		auth.POST("/check", a.check)
	}
}

func (a *Auth) signUp(c *gin.Context) {
	type RequestBody struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}

	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		a.ValidateError(c, err, http.StatusBadRequest)
		return
	}

	// check existancy
	var exists models.User
	if err := a.GetDB().Where("username = ?", body.Username).First(&exists).Error; err == nil {
		a.ValidateError(c, err, http.StatusConflict)
		return
	}

	hash, hashErr := models.Hash(body.Password)
	if hashErr != nil {
		a.ResponseError(c, errors.New("password hash error"))
		return
	}

	// create user
	user := models.User{
		Username:     body.Username,
		PasswordHash: hash,
	}

	a.GetDB().NewRecord(user)
	a.GetDB().Create(&user)

	serialized := user.Serialize()
	token, _ := models.GenerateToken(serialized)
	c.JSON(http.StatusOK, common.JSON{
		"user":  serialized,
		"token": token,
	})
}


func (a *Auth) login(c *gin.Context) {
	type RequestBody struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		a.ValidateError(c, err, http.StatusBadRequest)
		return
	}

	var user models.User
	if err := a.GetDB().Where("username = ?", body.Username).First(&user).Error; err != nil {
		a.ValidateError(c, err, http.StatusNotFound)
		return
	}


	if !models.CheckHash(body.Password, user.PasswordHash) {
		a.ValidateError(c, errors.New("un authorized"), http.StatusUnauthorized)
		return
	}

	serialized := user.Serialize()
	token, _ := models.GenerateToken(serialized)
	c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, common.JSON{
		"user":  serialized,
		"token": token,
	})
}

// check API will renew token when token life is less than 3 days, otherwise, return null for token
func (a *Auth) check(c *gin.Context) {
	userRaw, ok := c.Get("user")
	if !ok {
		a.ValidateError(c, errors.New("un authorized"), http.StatusUnauthorized)
		return
	}

	user := userRaw.(models.User)
	tokenExpire := int64(c.MustGet("token_expire").(float64))
	now := time.Now().Unix()
	diff := tokenExpire - now

	serialized := user.Serialize()


	if diff < 60*60*24*3 {
		// renew token
		token, _ := models.GenerateToken(serialized)
		c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)
		c.JSON(http.StatusOK, common.JSON{
			"token": token,
			"user":  serialized,
		})
		return
	}

	c.JSON(http.StatusOK, common.JSON{
		"token": nil,
		"user":  serialized,
	})
}

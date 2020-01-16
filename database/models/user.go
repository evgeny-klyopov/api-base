package models

import (
	"api-base/lib/common"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		gorm.Model
		Username     string `json:"host" gorm:"unique;not null"`
		PasswordHash string `json:"host2" gorm:"not null;"`
	}
)

// Serialize serializes user data
func (u *User) Serialize() common.JSON {
	return common.JSON{
		"id":           u.ID,
		"username":     u.Username,
	}
}

func (u *User) Read(m common.JSON) {
	u.ID = uint(m["id"].(float64))
	u.Username = m["username"].(string)
}



//`json:"host" validate:"required,min=5"`


func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(data common.JSON) (string, error) {
	//  token is valid for 7days
	date := time.Now().Add(time.Hour * 24 * 7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": data,
		"exp":  date.Unix(),
	})

	// get path from root dir
	pwd, _ := os.Getwd()
	keyPath := pwd + "/jwtsecret.key"

	key, readErr := ioutil.ReadFile(keyPath)
	if readErr != nil {
		return "", readErr
	}
	tokenString, err := token.SignedString(key)
	return tokenString, err
}


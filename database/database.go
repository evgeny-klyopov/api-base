package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // configures mysql driver
)

// Initialize initializes the database
func Initialize() (*gorm.DB, error) {
	dbConfig := os.Getenv("DB_CONFIG")
	dbLogMode, _ := strconv.ParseBool(os.Getenv("DB_LOG_MODE"))

	db, err := gorm.Open("mysql", dbConfig)
	if err != nil {
		panic(err)
	}

	db.LogMode(dbLogMode)

	fmt.Println("Connected to database")

	return db, err
}

package main

import (
	"api-base/api"
	"api-base/lib/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	// load .env environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}


	port := os.Getenv("PORT")
	app := gin.Default() // create gin app

	app.Use(middlewares.JWTMiddleware())
	api.NewApi(app).ApplyRoutes(nil)

	err = app.Run(":" + port)

	if err != nil {
		panic(err)
	}
}

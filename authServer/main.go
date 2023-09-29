package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var jwtKey string

func main() {
	InitDB()
	defer CloseDB()
	godotenv.Load()
	jwtKey = os.Getenv("JWT_KEY")

	r := gin.New()

	// Generate a new token
	r.POST("/token", func(c *gin.Context) {
		var loginData LoginData
		c.BindJSON(&loginData)
		createToken(loginData, c)
	})

	// Check if token is valid
	r.GET("/token", func(c *gin.Context) {
		isTokenValid(c)
	})

	// Get user uuid from token
	r.GET("/user", func(c *gin.Context) {
		returnTokenData(c)
	})

	fmt.Println("Auth-Server started. Port: 8082")
	r.Run(":8082")
}

package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type Player struct {
	ID       uuid.UUID `json:"id"`
	Color    string    `json:"color"`
	Username string    `json:"name"`
}

type Server struct {
	ID          uuid.UUID `json:"id"`
	PetName     string    `json:"petName"`
	Address     string    `json:"address"`
	PlayerCount int       `json:"playerCount"`
}

func main() {
	InitDB()
	defer CloseDB()

	r := gin.New()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"}
	r.Use(cors.New(config))

	r.GET("/listServer", func(c *gin.Context) {
		c.JSON(http.StatusOK, GetGameServers())
	})

	r.GET("/highscores", func(c *gin.Context) {
		functions := GetFunctions()
		c.IndentedJSON(http.StatusOK, functions)
	})

	fmt.Println("Lobby-Server started. Port: 8081")
	r.Run(":8081")
}

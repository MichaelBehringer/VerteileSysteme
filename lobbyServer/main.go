package main

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type Server struct {
	ID          uuid.UUID `json:"id"`
	Address     string    `json:"address"`
	PlayerCount int       `json:"playerCount"`
}

type Score struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Points   int       `json:"points"`
}

var servers = []Server{
	{ID: uuid.New(), Address: "ws://130.61.10.8:8080/ws", PlayerCount: 15},
	{ID: uuid.New(), Address: "ws://130.61.10.8:8081/ws -- geht nicht ", PlayerCount: 3},
	{ID: uuid.New(), Address: "ws://130.61.10.8:8082/ws -- geht nicht ", PlayerCount: 1},
}

var highscore = []Score{
	{ID: uuid.New(), Username: "Michael", Points: 99},
	{ID: uuid.New(), Username: "David", Points: 20},
	{ID: uuid.New(), Username: "Marco", Points: -10},
}

func main() {
	r := gin.New()

	r.GET("/listServer", func(c *gin.Context) {
		c.JSON(http.StatusOK, servers)
	})

	r.GET("/listScore", func(c *gin.Context) {
		c.JSON(http.StatusOK, highscore)
	})

	r.Run(":8090")
}

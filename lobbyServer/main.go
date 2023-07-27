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

var servers = []Server{
	{ID: uuid.New(), Address: "ws://130.61.10.8:8080/ws", PlayerCount: 15},
	{ID: uuid.New(), Address: "ws://130.61.10.8:8081/ws -- geht nicht ", PlayerCount: 3},
	{ID: uuid.New(), Address: "ws://130.61.10.8:8082/ws -- geht nicht ", PlayerCount: 1},
}

func main() {
	r := gin.New()

	r.GET("/listServer", func(c *gin.Context) {
		c.JSON(http.StatusOK, servers)
	})

	r.Run(":8090")
}

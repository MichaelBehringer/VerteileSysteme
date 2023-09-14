package main

import (
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

	r.GET("/player/:id", func(c *gin.Context) {
		var dummyPlayer = Player{ID: uuid.New(), Username: "Michael", Color: "blue"}
		c.JSON(http.StatusOK, dummyPlayer)
	})

	r.POST("/player", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/getUrl/:id", func(c *gin.Context) {
		id := c.Param("id")
		var port string
		ExecuteSQLRow("select g.Port from GameServer g where g.ID = ?", id).Scan(&port)

		c.JSON(http.StatusOK, port)
	})

	r.GET("/highscores", func(c *gin.Context) {
		functions := GetFunctions()
		c.IndentedJSON(http.StatusOK, functions)
	})

	r.Run("localhost:8090")
}

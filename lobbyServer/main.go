package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var colors []string

func main() {
	InitDB()
	defer CloseDB()

	// List of colors for random color new User
	colors = []string{
		"red", "green", "blue", "yellow", "maroon", "purple", "lime", "olive", "teal", "aqua",
		"orange", "pink", "brown", "gray", "beige", "fuchsia", "cyan", "magenta", "violet", "indigo",
		"navy", "silver", "gold", "hotPink", "turquoise", "lavender", "plum", "coral", "azure", "salmon",
	}

	r := gin.New()

	// Get all active game servers
	r.GET("/server", func(c *gin.Context) {
		c.JSON(http.StatusOK, getGameServers())
	})

	// Get highscore list of top Players
	r.GET("/highscore", func(c *gin.Context) {
		functions := getHighscore()
		c.IndentedJSON(http.StatusOK, functions)
	})

	// Get User Customaization of player in token
	r.GET("/user", func(c *gin.Context) {
		getUser(c)
	})

	// Update User Customaization of player in token
	r.POST("/user", func(c *gin.Context) {
		updateUser(c)
	})

	// Create a new user
	r.PUT("/user", func(c *gin.Context) {
		createUser(c)
	})

	fmt.Println("Lobby-Server started. Port: 8081")
	r.Run(":8081")
}

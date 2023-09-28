package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type Player struct {
	ID       uuid.UUID `json:"id"`
	Color    string    `json:"color"`
	Username string    `json:"name"`
}

type CustomPlayer struct {
	Skin     string `json:"skin"`
	Username string `json:"username"`
	Gamename string `json:"gamename"`
}

type Server struct {
	ID          uuid.UUID `json:"id"`
	PetName     string    `json:"petName"`
	Address     string    `json:"address"`
	PlayerCount int       `json:"playerCount"`
}

type AuthHeader struct {
	IDToken string `header:"Authorization"`
}

func ExtractToken(c *gin.Context) string {
	h := AuthHeader{}
	c.ShouldBindHeader(&h)
	idTokenHeader := strings.Split(h.IDToken, "Bearer ")
	if len(idTokenHeader) < 2 {
		return "noToken"
	}
	return idTokenHeader[1]
}

func GetUserId(c *gin.Context) string {
	token := ExtractToken(c)
	if token == "noToken" {
		return "error"
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://auth:8082/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, _ := client.Do(req)
	if res.StatusCode == 200 {
		body, _ := io.ReadAll(res.Body)
		var playerData map[string]interface{}
		json.Unmarshal(body, &playerData)
		userId, _ := playerData["userId"].(string)
		return userId
	}
	return "error"
}

func main() {
	InitDB()
	defer CloseDB()

	r := gin.New()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"}
	r.Use(cors.New(config))

	//TODO Endpoint schützen?
	r.GET("/server", func(c *gin.Context) {
		c.JSON(http.StatusOK, GetGameServers())
	})

	r.GET("/highscore", func(c *gin.Context) {
		functions := GetFunctions()
		c.IndentedJSON(http.StatusOK, functions)
	})

	r.GET("/user", func(c *gin.Context) {
		userId := GetUserId(c)
		if userId != "error" {
			var player CustomPlayer
			ExecuteSQLRow("SELECT Gamename, Skin, Username FROM Player WHERE ID=?", userId).Scan(&player.Gamename, &player.Skin, &player.Username)
			c.IndentedJSON(http.StatusOK, player)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	})

	r.POST("/user", func(c *gin.Context) {
		userId := GetUserId(c)
		if userId != "error" {
			var player CustomPlayer
			c.BindJSON(&player)
			ExecuteDDL("UPDATE Player SET Gamename = ?, Skin = ? where ID = ?", player.Gamename, player.Skin, userId)
			c.Status(http.StatusOK)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	})

	fmt.Println("Lobby-Server started. Port: 8081")
	r.Run(":8081")
}

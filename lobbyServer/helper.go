package main

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Return Highscore List from View
func getHighscore() []Highscore {
	results := ExecuteSQL("SELECT Username, Score FROM HighscoreList")
	highscores := []Highscore{}
	for results.Next() {
		var highscore Highscore
		results.Scan(&highscore.Name, &highscore.Highscore)
		highscores = append(highscores, highscore)
	}
	return highscores
}

// Return Active Game Servers from View
func getGameServers() []Server {
	results := ExecuteSQL("select ID, Servername, Servernumber, PlayerCounter from ActiveGameServer")
	servers := []Server{}
	for results.Next() {
		var server Server
		results.Scan(&server.ID, &server.PetName, &server.Address, &server.PlayerCount)
		servers = append(servers, server)
	}
	return servers
}

// get token string from header and check if it is valid
func extractToken(c *gin.Context) string {
	h := AuthHeader{}
	c.ShouldBindHeader(&h)
	idTokenHeader := strings.Split(h.IDToken, "Bearer ")
	if len(idTokenHeader) < 2 {
		return "noToken"
	}
	return idTokenHeader[1]
}

// Return userId from token; Auth-Service needed for getting Value of Token
func getUserId(c *gin.Context) string {
	token := extractToken(c)
	if token == "noToken" {
		return "error"
	}
	client := &http.Client{}
	// Adding Token to Header
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

// Create a new user
func createUser(c *gin.Context) {
	var player NewPlayer
	c.BindJSON(&player)
	var isAllowed bool
	// Check if username is already taken
	ExecuteSQLRow("SELECT COUNT(*) FROM Player WHERE USERNAME=?", player.Username).Scan(&isAllowed)
	if isAllowed {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Generate a hashed version of user password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(player.Password), bcrypt.DefaultCost)
	ExecuteDDL("INSERT INTO Player (ID, Username, Gamename, Skin, Passwort) VALUES(?, ?, ?, ?, ?)", uuid.New(), player.Username, player.Username, colors[rand.Intn(30)], hashedPassword)
	c.Status(http.StatusOK)
}

// Update User Customaization of player in token
func updateUser(c *gin.Context) {
	userId := getUserId(c)
	if userId != "error" {
		var player CustomPlayer
		c.BindJSON(&player)
		ExecuteDDL("UPDATE Player SET Gamename = ?, Skin = ? where ID = ?", player.Gamename, player.Skin, userId)
		c.Status(http.StatusOK)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

// Get User Customaization of player in token
func getUser(c *gin.Context) {
	userId := getUserId(c)
	if userId != "error" {
		var player CustomPlayer
		ExecuteSQLRow("SELECT Gamename, Skin, Username FROM Player WHERE ID=?", userId).Scan(&player.Gamename, &player.Skin, &player.Username)
		c.IndentedJSON(http.StatusOK, player)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

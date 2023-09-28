package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Player struct {
	Username string `json:"username"`
}

type ResponseText struct {
	Reason string `json:"reason"`
}

type AccessToken struct {
	AccessToken string `json:"accessToken"`
}

type AuthHeader struct {
	IDToken string `header:"Authorization"`
}

func CreateToken(loginData LoginData, c *gin.Context) {
	var (
		key         []byte
		token       *jwt.Token
		signedToken string
	)

	var isAllowed bool
	ExecuteSQLRow("SELECT COUNT(*) FROM Player WHERE USERNAME=? AND PASSWORT=?", loginData.Username, loginData.Password).Scan(&isAllowed)
	if !isAllowed {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	key = []byte("my_secret_key")
	token = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user":         loginData.Username,
			"creationTime": time.Now().UnixNano(),
		})
	signedToken, _ = token.SignedString(key)

	c.IndentedJSON(http.StatusOK, AccessToken{AccessToken: signedToken})
}

func ExtractToken(c *gin.Context) (bool, jwt.MapClaims) {
	h := AuthHeader{}
	c.ShouldBindHeader(&h)
	idTokenHeader := strings.Split(h.IDToken, "Bearer ")
	if len(idTokenHeader) < 2 {
		return false, nil
	}
	return parseToken(idTokenHeader[1])
}

func parseToken(tokenStr string) (bool, jwt.MapClaims) {
	claims := jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("my_secret_key"), nil
	})
	return (err == nil && tkn.Valid), claims
}

func isTokenValid(c *gin.Context) {
	isAllowed, claims := ExtractToken(c)
	if claims == nil || !isAllowed {
		c.AbortWithStatus(http.StatusUnauthorized)
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func ReturnUserName(c *gin.Context) {
	isAllowed, claims := ExtractToken(c)
	fmt.Println(isAllowed, claims)
	if !isAllowed {
		c.AbortWithStatus(http.StatusUnauthorized)
	} else {
		username, _ := claims["user"].(string)
		player := Player{Username: username}
		c.IndentedJSON(http.StatusOK, player)
	}
}

func main() {
	InitDB()
	defer CloseDB()

	r := gin.New()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"}
	r.Use(cors.New(config))

	r.POST("/token", func(c *gin.Context) {
		var loginData LoginData
		c.BindJSON(&loginData)
		CreateToken(loginData, c)
	})

	r.GET("/token", func(c *gin.Context) {
		isTokenValid(c)
	})

	r.GET("/user", func(c *gin.Context) {
		ReturnUserName(c)
	})

	r.Run(":8081")
}

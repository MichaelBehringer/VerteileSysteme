package main

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Generate a new token
func createToken(loginData LoginData, c *gin.Context) {
	var (
		key         []byte
		token       *jwt.Token
		signedToken string
	)

	var password string
	ExecuteSQLRow("SELECT PASSWORT FROM Player WHERE USERNAME=?", loginData.Username).Scan(&password)
	if password == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// compare password with hash
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(loginData.Password))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var userId string
	ExecuteSQLRow("SELECT ID FROM Player WHERE USERNAME=?", loginData.Username).Scan(&userId)

	key = []byte(jwtKey)
	// create token with userUUID and current Time
	token = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user":         userId,
			"creationTime": time.Now().UnixNano(),
		})
	signedToken, _ = token.SignedString(key)

	c.IndentedJSON(http.StatusOK, AccessToken{AccessToken: signedToken})
}

// get token string from header and check if it is valid
func extractToken(c *gin.Context) (bool, jwt.MapClaims) {
	h := AuthHeader{}
	c.ShouldBindHeader(&h)
	idTokenHeader := strings.Split(h.IDToken, "Bearer ")
	if len(idTokenHeader) < 2 {
		return false, nil
	}
	return parseToken(idTokenHeader[1])
}

// get token claims from token string
func parseToken(tokenStr string) (bool, jwt.MapClaims) {
	claims := jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	return (err == nil && tkn.Valid), claims
}

// check if token is valid
func isTokenValid(c *gin.Context) {
	isAllowed, claims := extractToken(c)
	if claims == nil || !isAllowed {
		c.AbortWithStatus(http.StatusUnauthorized)
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

// Get user uuid from token
func returnTokenData(c *gin.Context) {
	isAllowed, claims := extractToken(c)
	if !isAllowed {
		c.AbortWithStatus(http.StatusUnauthorized)
	} else {
		userId, _ := claims["user"].(string)
		player := Player{UserId: userId}
		c.IndentedJSON(http.StatusOK, player)
	}
}

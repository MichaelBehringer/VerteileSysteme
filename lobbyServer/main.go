package main

import (
	"net/http"
	"time"

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
	LastCall    time.Time `json:"lastCall"`
}

type Score struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Points   int       `json:"points"`
}

var servers = []Server{}

var highscore = []Score{
	{ID: uuid.New(), Username: "Michael", Points: 99},
	{ID: uuid.New(), Username: "David", Points: 20},
	{ID: uuid.New(), Username: "Marco", Points: -10},
}

func addOrUpdateServer(newServer Server) {
	for i, existingServer := range servers {
		if existingServer.ID == newServer.ID {
			servers[i] = newServer
			return
		}
	}

	// ID nicht gefunden, hinzufügen
	servers = append(servers, newServer)
}

func areLastCallsDifferent(server1 Server, requestTime time.Time, threshold time.Duration) bool {
	timeDifference := server1.LastCall.Sub(requestTime)
	return timeDifference > -threshold
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
		requestTime := time.Now()
		var returnServer = []Server{}
		for _, existingServer := range servers {
			if areLastCallsDifferent(existingServer, requestTime, 15*time.Second) {
				returnServer = append(returnServer, existingServer)
			}
		}
		c.JSON(http.StatusOK, returnServer)
	})

	r.GET("/listScore", func(c *gin.Context) {
		c.JSON(http.StatusOK, highscore)
	})

	r.GET("/player/:id", func(c *gin.Context) {
		var dummyPlayer = Player{ID: uuid.New(), Username: "Michael", Color: "blue"}
		c.JSON(http.StatusOK, dummyPlayer)
	})

	r.POST("/player", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/getUrl/:id", func(c *gin.Context) {
		url := "invalidID"
		id := c.Param("id")
		for _, v := range servers {
			if v.ID.String() == id {
				url = v.Address
			}
		}

		c.JSON(http.StatusOK, url)
	})

	r.GET("/highscores", func(c *gin.Context) {
		// functions := GetFunctions()
		var functions = []Highscore{{Highscore: 69, Name: "Franzel"}}
		c.IndentedJSON(http.StatusOK, functions)
	})

	r.POST("/registerGameServer", func(c *gin.Context) {
		var requestData map[string]interface{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		requestTime := time.Now()
		serverId, _ := uuid.Parse(requestData["id"].(string))
		petName := requestData["petName"].(string)
		serverAddress := requestData["address"].(string)
		serverPlayerCount := int(requestData["playerCounter"].(float64))
		newServer := Server{ID: serverId, PetName: petName, Address: serverAddress, PlayerCount: serverPlayerCount, LastCall: requestTime}
		addOrUpdateServer(newServer)

		c.Status(http.StatusOK)
	})

	r.Run("localhost:8090")
}

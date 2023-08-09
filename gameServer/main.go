package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type myConnection struct {
	Key  uuid.UUID       `json:"key"`
	Conn *websocket.Conn `json:"connection"`
}

type myKoordinates struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Color int `json:"color"`
}

var listConnections []myConnection
var listKoordinates map[uuid.UUID]myKoordinates
var myMutex sync.RWMutex

var upgrader = websocket.Upgrader{
	Subprotocols: []string{"demo-chat"},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Fehler beim Upgrade der Verbindung:", err)
		return
	}
	defer conn.Close()

	randomColor := rand.Intn(10)
	id := uuid.New()
	s := myConnection{Key: id, Conn: conn}
	listConnections = append(listConnections, s)
	defer removeUnusedConnection(id)

	fmt.Printf("New User: %s\n", id)

	fmt.Println("WebSocket-Verbindung hergestellt")

	for {
		// Nachrichten vom Client empfangen
		messageType, message, _ := conn.ReadMessage()
		if messageType == -1 {
			break
		}
		koordinates := myKoordinates{}
		json.Unmarshal([]byte(message), &koordinates)
		koordinates.Color = randomColor
		myMutex.Lock()
		listKoordinates[id] = koordinates
		myMutex.Unlock()
	}
}

func removeUnusedConnection(key uuid.UUID) {
	fmt.Printf("removing one: %s\n", key)

	for idx, v := range listConnections {
		if v.Key == key {
			listConnections = remove(listConnections, idx)
		}
	}
	delete(listKoordinates, key)

}

func remove(s []myConnection, i int) []myConnection {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func main() {
	rand.Seed(time.Now().UnixNano())
	listKoordinates = make(map[uuid.UUID]myKoordinates)
	myMutex = sync.RWMutex{}
	ticker := time.NewTicker(time.Second / 100)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				// fmt.Printf("thread")
				j, _ := json.Marshal(listKoordinates)
				for _, singleConn := range listConnections {
					singleConn.Conn.WriteMessage(1, j)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	r := gin.New()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"}

	r.Use(cors.New(config))

	r.GET("/ws", func(c *gin.Context) {
		handleWebSocketConnection(c.Writer, c.Request)
	})

	r.GET("/players", func(c *gin.Context) {
		c.JSON(http.StatusOK, listConnections)
	})

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, listKoordinates)
	})

	fmt.Println("WebSocket-Server gestartet. Lausche auf http://localhost:8080/ws")

	r.Run(":8080")
}

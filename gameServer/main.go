package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type connectionObj struct {
	Key  uuid.UUID       `json:"key"`
	Conn *websocket.Conn `json:"connection"`
}

type playerObj struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Color int     `json:"c"`
	Size  int     `json:"s"`
}

type npcObj struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Color int `json:"c"`
}

type transfairObj struct {
	Player map[uuid.UUID]playerObj `json:"player"`
	Npc    map[uuid.UUID]npcObj    `json:"npc"`
}

var listConnections []connectionObj
var listPlayerKoordinates map[uuid.UUID]playerObj
var listNpcKoordinates map[uuid.UUID]npcObj
var myMutex sync.RWMutex
var myNpcMutex sync.RWMutex

func calcNewPoint(xStart float64, yStart float64, xEnd float64, yEnd float64) (float64, float64) {
	vectorX := xEnd - xStart
	vectorY := yEnd - yStart

	lenghtVector := math.Sqrt(math.Pow(vectorX, 2) + math.Pow(vectorY, 2))
	if lenghtVector < 1 {
		return xStart, yStart
	}

	normalizedX := vectorX / lenghtVector
	normalizedY := vectorY / lenghtVector

	stepSize := 2

	newX := xStart + normalizedX*float64(stepSize)
	newY := yStart + normalizedY*float64(stepSize)

	return newX, newY
}

// copied from stackoverflow, converts websocket
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

	// color for player
	randomColor := rand.Intn(9)
	// color for player
	size := 20
	// id for player
	id := uuid.New()
	// create player object and add to connectionList; delete from list on function end
	s := connectionObj{Key: id, Conn: conn}
	listConnections = append(listConnections, s)
	defer removeUnusedConnection(id)

	listPlayerKoordinates[id] = playerObj{X: float64(rand.Intn(1000)), Y: float64(rand.Intn(700)), Color: randomColor, Size: 20}

	// debug
	fmt.Printf("New User: %s\n", id)
	fmt.Println("WebSocket-Verbindung hergestellt")

	for {
		// Message recive
		messageType, message, _ := conn.ReadMessage()
		// messageType == -1 connection was closed
		if messageType == -1 {
			break
		}

		// create koordinates obj; parse recived koordinates to obj; mutex needed because of multithreding; update koordinates in hashmap
		koordinates := playerObj{}
		json.Unmarshal([]byte(message), &koordinates)
		// koordinates.Color = randomColor
		myMutex.Lock()
		oldPlayerKoords := listPlayerKoordinates[id]
		newX, newY := calcNewPoint(oldPlayerKoords.X, oldPlayerKoords.Y, koordinates.X, koordinates.Y)
		oldPlayerKoords.X = newX
		oldPlayerKoords.Y = newY
		oldPlayerKoords.Size = size
		listPlayerKoordinates[id] = oldPlayerKoords
		myMutex.Unlock()

		myNpcMutex.Lock()
		for key, value := range listNpcKoordinates {
			if (float64(value.X+size/2) >= oldPlayerKoords.X && float64(value.X-size/2) <= oldPlayerKoords.X) && (float64(value.Y+size/2) >= oldPlayerKoords.Y && float64(value.Y-size/2) <= oldPlayerKoords.Y) {
				delete(listNpcKoordinates, key)
				size++
			}
		}
		myNpcMutex.Unlock()

		myMutex.Lock()
		for key, value := range listPlayerKoordinates {
			if key != id && (value.X+float64((size+value.Size)/2) >= oldPlayerKoords.X && value.X-float64((size+value.Size)/2) <= oldPlayerKoords.X) && (value.Y+float64((size+value.Size)/2) >= oldPlayerKoords.Y && value.Y-float64((size+value.Size)/2) <= oldPlayerKoords.Y) {
				if size > value.Size {
					listPlayerKoordinates[key] = playerObj{X: float64(rand.Intn(1000)), Y: float64(rand.Intn(700)), Color: value.Color, Size: 20}
				} else {
					listPlayerKoordinates[id] = playerObj{X: float64(rand.Intn(1000)), Y: float64(rand.Intn(700)), Color: randomColor, Size: 20}
				}
			}
		}
		myMutex.Unlock()
	}
}

// removes connection from slice
func removeUnusedConnection(key uuid.UUID) {
	fmt.Printf("removing one: %s\n", key)

	for idx, v := range listConnections {
		if v.Key == key {
			listConnections = remove(listConnections, idx)
		}
	}
	delete(listPlayerKoordinates, key)

}

func remove(s []connectionObj, i int) []connectionObj {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func main() {
	// srand
	rand.Seed(time.Now().UnixNano())
	// initialize hashmap and mutex
	listPlayerKoordinates = make(map[uuid.UUID]playerObj)
	listNpcKoordinates = make(map[uuid.UUID]npcObj)

	playerUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	listPlayerKoordinates[playerUUID] = playerObj{X: -1000, Y: -1000, Color: 0, Size: 1}

	npcUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
	listNpcKoordinates[npcUUID] = npcObj{X: -1000, Y: -1000, Color: 0}

	// thread for broadcast
	myMutex = sync.RWMutex{}
	myNpcMutex = sync.RWMutex{}
	ticker := time.NewTicker(time.Second / 100)
	ticker4 := time.NewTicker(time.Second / 10)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				// fmt.Printf("thread")
				// myMutex.Lock()
				// s := transfairObj{Player: listPlayerKoordinates, Npc: listNpcKoordinates}
				// myMutex.Unlock()
				// j, _ := json.Marshal(listPlayerKoordinates)
				myMutex.Lock()
				for _, singleConn := range listConnections {
					// singleConn.Conn.WriteMessage(1, j)
					singleConn.Conn.WriteJSON(listPlayerKoordinates)
				}
				myMutex.Unlock()
			case <-ticker4.C:
				myNpcMutex.Lock()
				for _, singleConn := range listConnections {
					singleConn.Conn.WriteJSON(listNpcKoordinates)
				}
				myNpcMutex.Unlock()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	ticker2 := time.NewTicker(time.Second)
	quit2 := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker2.C:
				for {
					if len(listNpcKoordinates) >= 100 {
						break
					}
					myNpcMutex.Lock()
					npc := npcObj{X: rand.Intn(1000), Y: rand.Intn(700), Color: rand.Intn(10)}
					listNpcKoordinates[uuid.New()] = npc
					myNpcMutex.Unlock()
				}
				// j, _ := json.Marshal(listNpcKoordinates)
				// for _, singleConn := range listConnections {
				// 	singleConn.Conn.WriteMessage(1, j)
				// }
			case <-quit2:
				ticker2.Stop()
				return
			}
		}
	}()

	r := gin.New()

	// corse error; maybe delete later ?!?
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"}

	r.Use(cors.New(config))

	// main websocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		handleWebSocketConnection(c.Writer, c.Request)
	})

	// debug enpoints
	r.GET("/players", func(c *gin.Context) {
		c.JSON(http.StatusOK, listConnections)
	})

	r.GET("/test", func(c *gin.Context) {
		s := transfairObj{Player: listPlayerKoordinates, Npc: listNpcKoordinates}
		c.JSON(http.StatusOK, s)
	})

	fmt.Println("WebSocket-Server gestartet. Lausche auf http://localhost:8080/ws")

	r.Run(":8080")
}

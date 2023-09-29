package main

import (
	"fmt"
	"net/http"

	"github.com/dhconnelly/rtreego"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

var mapBoundary = 5000.0

// List of all connections
var listConnections []connectionObj

// List of all players; mapIdToPlayer: uuid to index of listPlayerKoordinates
var listPlayerKoordinates [30]gameObj
var mapIdToPlayer map[uuid.UUID]int

// List of the Location where each Player wants to travel
var arrPlayerTarget [30]targetObj

// List of all NPCs
var listNpcKoordinates [1200]gameObj

// Current Chat Messages to send to the clients
var listMessages []transfairChatMessage

// Game server variables
var playerCounter int
var gameServerId uuid.UUID

var stack *Stack
var colors []string

// global Tree for NPC, n log n search time
var treeNpc *rtreego.Rtree

// Websocket connection object
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func main() {
	InitDB()
	defer CloseDB()
	treeNpc = rtreego.NewTree(2, 25, 50)
	mapIdToPlayer = make(map[uuid.UUID]int)
	stack = newStack()
	playerCounter = 0
	gameServerId = uuid.New()

	colors = []string{
		"red", "green", "blue", "yellow", "maroon", "purple", "lime", "olive", "teal", "aqua",
		"orange", "pink", "brown", "gray", "beige", "fuchsia", "cyan", "magenta", "violet", "indigo",
		"navy", "silver", "gold", "hotPink", "turquoise", "lavender", "plum", "coral", "azure", "salmon",
	}

	r := gin.New()

	// main websocket endpoint for game communication
	r.GET("/ws/:token", func(c *gin.Context) {
		token := c.Param("token")
		handleWebSocketConnection(c.Writer, c.Request, token)
	})

	// initialisation of the game
	initNPCs()
	initStack()

	// starting the go routines (like threads) for game checks
	go movePlayer()
	go checkCollission()
	go sendUpdate()

	// generating name, containerNo and start live checks and highscore db updates
	serverPetName := getPetName()
	containerNo := getContainerNo()
	go gameServerAlive(containerNo, serverPetName)
	go updateHighscore(containerNo)

	fmt.Println("Game-Server started. Port: 8080")
	r.Run(":8080")
}

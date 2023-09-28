package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/dhconnelly/rtreego"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type connectionObj struct {
	Key       uuid.UUID       `json:"key"`
	Conn      *websocket.Conn `json:"connection"`
	ConnMutex *sync.Mutex     `json:"connectionMutex"`
}

type gameObj struct {
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Color      string  `json:"color"`
	Name       string  `json:"name"`
	Size       float64 `json:"size"`
	IsLoggedIn bool    `json:"isLoggendIn"`
}

type targetObj struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type finalTransfairObj struct {
	Player      transfairPlayer   `json:"player"`
	OtherPlayer []transfairPlayer `json:"otherPlayer"`
	NPC         []transfairNpc    `json:"npc"`
	Score       []transfairScore  `json:"score"`
}

type transfairPlayer struct {
	Id    uuid.UUID `json:"id"`
	Color string    `json:"color"`
	Name  string    `json:"name"`
	Size  float64   `json:"size"`
	X     float64   `json:"x"`
	Y     float64   `json:"y"`
}

type transfairNpc struct {
	Color string  `json:"color"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
}

type transfairScore struct {
	Name  string  `json:"name"`
	Score float64 `json:"highscore"`
}

type AuthHeader struct {
	IDToken string `header:"Authorization"`
}

var mapBoundary = 5000.0

var listConnections []connectionObj
var listPlayerKoordinates [30]gameObj
var listNpcKoordinates [1200]gameObj
var arrPlayerTarget [30]targetObj
var mapIdToPlayer map[uuid.UUID]int

var stack *Stack
var playerCounter int
var gameServerId uuid.UUID
var colors []string

var treeNpc *rtreego.Rtree

var words = flag.Int("words", 2, "The number of words in the pet name")
var separator = flag.String("separator", "-", "The separator between words in the pet name")

func calcNewPoint(xStart float64, yStart float64, xEnd float64, yEnd float64, size float64) (float64, float64) {
	//map eingrenzen auf 0 bis 5k
	vectorX := xEnd - xStart
	vectorY := yEnd - yStart

	lenghtVector := math.Sqrt(math.Pow(vectorX, 2) + math.Pow(vectorY, 2))
	if lenghtVector == 0 {
		return xStart, yStart
	}

	normalizedX := vectorX / lenghtVector
	normalizedY := vectorY / lenghtVector

	stepSize := 2.0 + 5.0/size

	if lenghtVector < 100 {
		stepSize = stepSize * lenghtVector * 0.01
	}

	newX := xStart + normalizedX*stepSize
	newY := yStart + normalizedY*stepSize

	if newX > mapBoundary {
		newX = mapBoundary
	} else if newX < 0.0 {
		newX = 0.0
	}
	if newY > mapBoundary {
		newY = mapBoundary
	} else if newY < 0.0 {
		newY = 0.0
	}

	return newX, newY
}

// copied from stackoverflow, converts websocket
var upgrader = websocket.Upgrader{
	Subprotocols: []string{"demo-chat"},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request, token string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Fehler beim Upgrade der Verbindung:", err)
		return
	}
	defer conn.Close()

	if playerCounter > 30 {
		fmt.Println("Max Player reached")
		return
	}

	var username string
	var skin string
	var uuidNode uuid.UUID
	isLoggedIn := token != "undefined"
	if isLoggedIn {
		username, skin, uuidNode = getTokenData(token)
	} else {
		username, skin, uuidNode = getRandomTokenData()
	}

	playerCounter += 1
	playerId, _ := stack.Pop()
	defer stack.Push(playerId)

	mapIdToPlayer[uuidNode] = playerId

	// create player object and add to connectionList; delete from list on function end
	s := connectionObj{Key: uuidNode, Conn: conn, ConnMutex: &sync.Mutex{}}
	listConnections = append(listConnections, s)
	defer removeUnusedConnection(uuidNode)

	listPlayerKoordinates[playerId] = gameObj{X: randFloat(0, mapBoundary), Y: randFloat(0, mapBoundary), Color: skin, Size: 20, Name: username, IsLoggedIn: isLoggedIn}

	// debug
	fmt.Printf("New User: %d:-- %s\n", playerId, uuidNode)
	fmt.Println("WebSocket-Verbindung hergestellt")

	for {
		messageType, message, _ := conn.ReadMessage()
		if messageType == -1 {
			break
		}

		koordinates := targetObj{}
		json.Unmarshal([]byte(message), &koordinates)
		arrPlayerTarget[playerId] = koordinates
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
	// delete(listPlayerKoordinates, key)
	delete(mapIdToPlayer, key)
	playerCounter -= 1

}

func remove(s []connectionObj, i int) []connectionObj {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func getPetName() string {
	flag.Parse()
	return petname.Generate(*words, *separator)
}

func getPetNameSingle() string {
	return petname.Generate(1, "-")
}

func updateHighscore(containerNo string) {
	for range time.Tick(time.Second * 15) {
		for idx, value := range mapIdToPlayer {
			player := listPlayerKoordinates[value]
			if player.IsLoggedIn {
				ExecuteDDL("CALL InsertUpdateHighscore(?, ?, ?)", gameServerId, idx, player.Size)
			}
		}
	}
}

func gameServerAlive(containerNo string, petName string) {
	execGameServerAlive(containerNo, petName)
	for range time.Tick(time.Second * 10) {
		execGameServerAlive(containerNo, petName)
	}
}

func execGameServerAlive(containerNo string, petName string) {
	ExecuteDDL("CALL InsertUpdateGameServer(?, ?, ?, ?)", gameServerId, petName, containerNo, playerCounter)
}

func initNPCs() {
	for i := 0; i < 1200; i++ {
		npc := gameObj{X: randFloat(0, mapBoundary), Y: randFloat(0, mapBoundary), Color: colors[rand.Intn(30)]}
		listNpcKoordinates[i] = npc
		treeNpc.Insert(Circle{Id: 0, X: npc.X, Y: npc.Y, Radius: 10, Color: npc.Color})
	}
}

func initStack() {
	for i := 0; i < 30; i++ {
		stack.Push(i)
	}
}

func movePlayer() {
	for range time.Tick(time.Second / 100) {
		for _, value := range mapIdToPlayer {
			oldPlayerKoords := listPlayerKoordinates[value]
			playerTarget := arrPlayerTarget[value]
			newX, newY := calcNewPoint(oldPlayerKoords.X, oldPlayerKoords.Y, playerTarget.X+oldPlayerKoords.X, playerTarget.Y+oldPlayerKoords.Y, oldPlayerKoords.Size)
			listPlayerKoordinates[value].X = newX
			listPlayerKoordinates[value].Y = newY
		}
	}
}

func checkCollission() {
	for range time.Tick(time.Second / 10) {
		npcCollision()
		playerCollision()
	}
}

func sendUpdate() {
	count := 0
	for range time.Tick(time.Second / 100) {
		count++
		sendNpcUpdate := count%25 == 0
		sendScoreUpdate := count%200 == 0
		for _, singleConn := range listConnections {
			go sendSingleUpdate(singleConn, sendNpcUpdate, sendScoreUpdate)
		}
	}
}

func sendSingleUpdate(singleConn connectionObj, sendNpcUpdate bool, sendScoreUpdate bool) {
	var objPlayerT transfairPlayer
	var objOtherPlayerT []transfairPlayer
	for key, value := range mapIdToPlayer {
		player := listPlayerKoordinates[value]
		transfiarPlayerObj := transfairPlayer{Id: key, Color: player.Color, Name: player.Name, Size: player.Size, X: player.X, Y: player.Y}
		if key == singleConn.Key {
			objPlayerT = transfiarPlayerObj
		} else {
			objOtherPlayerT = append(objOtherPlayerT, transfiarPlayerObj)
		}
	}
	var objTransfairNpc []transfairNpc
	if sendNpcUpdate {
		objTransfairNpc = visibleNPC(objPlayerT)
	}
	var objTransfairScore []transfairScore
	if sendScoreUpdate {
		objTransfairScore = getScore(objPlayerT)
	}
	finalObjT := finalTransfairObj{Player: objPlayerT, OtherPlayer: objOtherPlayerT, NPC: objTransfairNpc, Score: objTransfairScore}
	singleConn.ConnMutex.Lock()
	singleConn.Conn.WriteJSON(finalObjT)
	singleConn.ConnMutex.Unlock()
}

func main() {
	InitDB()
	defer CloseDB()
	treeNpc = rtreego.NewTree(2, 25, 50)
	mapIdToPlayer = make(map[uuid.UUID]int)
	stack = NewStack()
	playerCounter = 0
	gameServerId = uuid.New()

	colors = []string{
		"red", "green", "blue", "yellow", "maroon", "purple", "lime", "olive", "teal", "aqua",
		"orange", "pink", "brown", "gray", "beige", "fuchsia", "cyan", "magenta", "violet", "indigo",
		"navy", "silver", "gold", "hotPink", "turquoise", "lavender", "plum", "coral", "azure", "salmon",
	}

	r := gin.New()

	// corse error; maybe delete later ?!?
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"}

	r.Use(cors.New(config))

	// main websocket endpoint
	r.GET("/ws/:token", func(c *gin.Context) {
		token := c.Param("token")
		handleWebSocketConnection(c.Writer, c.Request, token)
	})

	initNPCs()
	initStack()

	go movePlayer()
	go checkCollission()
	go sendUpdate()

	serverPetName := getPetName()
	containerNo := getContainerNo()
	go gameServerAlive(containerNo, serverPetName)
	go updateHighscore(containerNo)

	fmt.Println("Game-Server started. Port: 8080")
	r.Run(":8080")
}

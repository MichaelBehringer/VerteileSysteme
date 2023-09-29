package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

// Recives new Websocket connections, handles messages from Client
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

	// Setting up player
	var username string
	var skin string
	var uuidNode uuid.UUID
	isLoggedIn := token != "undefined"
	if isLoggedIn {
		username, skin, uuidNode = getTokenData(token)
	} else {
		username, skin, uuidNode = getRandomTokenData()
	}

	if mapIdToPlayer[uuidNode] != 0 {
		fmt.Println("Player already in game")
		return
	}

	playerCounter += 1
	playerId, _ := stack.pop()
	defer stack.push(playerId)

	mapIdToPlayer[uuidNode] = playerId

	// create player object and add to connectionList; delete from list on function end
	s := connectionObj{Key: uuidNode, Conn: conn, ConnMutex: &sync.Mutex{}}
	listConnections = append(listConnections, s)
	defer removeUnusedConnection(uuidNode)

	listPlayerKoordinates[playerId] = gameObj{X: randFloat(0, mapBoundary), Y: randFloat(0, mapBoundary), Color: skin, Size: 20, Name: username, IsLoggedIn: isLoggedIn}

	fmt.Printf("New User: %d:-- %s\n", playerId, uuidNode)

	for {
		messageType, message, _ := conn.ReadMessage()
		if messageType == -1 {
			// Client disconnected
			break
		}
		socketMode := socketMode{}
		json.Unmarshal([]byte(message), &socketMode)
		if socketMode.Mode == "pos" {
			// Client sends new coordinates
			koordinates := targetObj{}
			json.Unmarshal([]byte(message), &koordinates)
			arrPlayerTarget[playerId] = koordinates
		} else {
			// Client sends chat message
			socketMessage := socketMessage{}
			json.Unmarshal([]byte(message), &socketMessage)
			msg := transfairChatMessage{Name: username, Message: socketMessage.Message, Size: int(listPlayerKoordinates[playerId].Size)}
			listMessages = append(listMessages, msg)
		}
	}
}

// Get the container replica number
// User Docker client to get name of the container
// Needed so that the container can be identified in the database
func getContainerNo() string {
	cli, _ := client.NewClientWithOpts(client.FromEnv)
	containerID, _ := os.Hostname()
	containerInfo, _ := cli.ContainerInspect(context.Background(), containerID)
	containerName := containerInfo.Name[1:]
	splitFunc := func(r rune) bool {
		return r == '-' || r == '_'
	}
	containerName = strings.FieldsFunc(containerName, splitFunc)[2]
	return containerName
}

// Sending Updates to Player 10ms; Messages are Customized to the Player
func sendUpdate() {
	count := 0
	for range time.Tick(time.Second / 100) {
		// Performance Tuning; dont send npc Location and Score every message
		count++
		sendNpcUpdate := count%25 == 0
		sendScoreUpdate := count%200 == 0

		var objTransfairMessage []transfairChatMessage
		if sendNpcUpdate && len(listMessages) > 0 {
			objTransfairMessage = make([]transfairChatMessage, len(listMessages))
			copy(objTransfairMessage, listMessages)
			listMessages = []transfairChatMessage{}
		}

		var objTransfairScore []transfairScore
		if sendScoreUpdate {
			objTransfairScore = getScore()
		}

		for _, singleConn := range listConnections {
			go sendSingleUpdate(singleConn, sendNpcUpdate, objTransfairMessage, objTransfairScore)
		}
	}
}

func sendSingleUpdate(singleConn connectionObj, sendNpcUpdate bool, objTransfairMessage []transfairChatMessage, objTransfairScore []transfairScore) {
	var objPlayerT transfairPlayer
	var objOtherPlayerT []transfairPlayer

	// Current Player
	userPlayer := listPlayerKoordinates[mapIdToPlayer[singleConn.Key]]
	objPlayerT = transfairPlayer{Color: userPlayer.Color, Name: userPlayer.Name, Size: userPlayer.Size, X: userPlayer.X, Y: userPlayer.Y}

	for key, value := range mapIdToPlayer {
		// Other Players, but only if they are in a 1200px radius
		player := listPlayerKoordinates[value]
		if key != singleConn.Key && math.Abs(userPlayer.X-player.X) < 1200 && math.Abs(userPlayer.Y-player.Y) < 1200 {
			transfairPlayerObj := transfairPlayer{Color: player.Color, Name: player.Name, Size: player.Size, X: player.X, Y: player.Y}
			objOtherPlayerT = append(objOtherPlayerT, transfairPlayerObj)
		}
	}

	// NPCs in a 500px radius
	var objTransfairNpc []transfairNpc
	if sendNpcUpdate {
		objTransfairNpc = visibleNPC(objPlayerT)
	}
	// Send Update to Client; Mutex needed because of concurrency
	finalObjT := finalTransfairObj{Player: objPlayerT, OtherPlayer: objOtherPlayerT, NPC: objTransfairNpc, Score: objTransfairScore, Message: objTransfairMessage}
	singleConn.ConnMutex.Lock()
	singleConn.Conn.WriteJSON(finalObjT)
	singleConn.ConnMutex.Unlock()
}

// Heartbeat to database every 10s
func gameServerAlive(containerNo string, petName string) {
	execGameServerAlive(containerNo, petName)
	for range time.Tick(time.Second * 10) {
		execGameServerAlive(containerNo, petName)
	}
}

// calls stored procedure for updating gameServers; Logic when Game Server is updated is in the Database
func execGameServerAlive(containerNo string, petName string) {
	ExecuteDDL("CALL InsertUpdateGameServer(?, ?, ?, ?)", gameServerId, petName, containerNo, playerCounter)
}

func removeConnection(s []connectionObj, i int) []connectionObj {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// removes connection from slice
func removeUnusedConnection(key uuid.UUID) {
	fmt.Printf("removing one: %s\n", key)

	for idx, v := range listConnections {
		if v.Key == key {
			listConnections = removeConnection(listConnections, idx)
		}
	}
	delete(mapIdToPlayer, key)
	playerCounter -= 1

}

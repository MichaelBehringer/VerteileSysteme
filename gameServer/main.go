package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type myConnection struct {
	Key  uuid.UUID       `json:"key"`
	Conn *websocket.Conn `json:"connection"`
}

var listConnections []myConnection

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

		fmt.Printf("Nachricht erhalten: %d %s\n", messageType, message)

		for _, singleConn := range listConnections {
			fmt.Printf("sendin to: %s\n", singleConn.Key)
			singleConn.Conn.WriteMessage(1, []byte("gele Projektgruppe: "+string(message)))
		}
	}
}

func removeUnusedConnection(key uuid.UUID) {
	fmt.Printf("removing one: %s\n", key)

	for idx, v := range listConnections {
		if v.Key == key {
			listConnections = remove(listConnections, idx)
		}
	}

}

func remove(s []myConnection, i int) []myConnection {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func main() {
	r := gin.New()
	r.GET("/ws", func(c *gin.Context) {
		handleWebSocketConnection(c.Writer, c.Request)
	})

	r.GET("/players", func(c *gin.Context) {
		c.JSON(http.StatusOK, listConnections)
	})

	fmt.Println("WebSocket-Server gestartet. Lausche auf http://localhost:8080/ws")

	r.Run(":8080")
}

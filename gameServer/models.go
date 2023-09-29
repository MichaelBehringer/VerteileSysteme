package main

import (
	"sync"

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

type socketMode struct {
	Mode string `json:"mode"`
}

type socketMessage struct {
	Message string `json:"message"`
}

type Circle struct {
	Id          int
	Name, Color string
	X, Y        float64
	Radius      float64
}

type Player struct {
	Gamename string `json:"gamename"`
	Skin     string `json:"skin"`
}

type Stack struct {
	items []int
}

type AuthHeader struct {
	IDToken string `header:"Authorization"`
}

type AccessToken struct {
	AccessToken string `json:"accessToken"`
}

//-------------------------------------------------------
// Transfair Objects

type finalTransfairObj struct {
	Player      transfairPlayer        `json:"player"`
	OtherPlayer []transfairPlayer      `json:"otherPlayer"`
	NPC         []transfairNpc         `json:"npc"`
	Score       []transfairScore       `json:"score"`
	Message     []transfairChatMessage `json:"message"`
}

type transfairPlayer struct {
	Color string  `json:"color"`
	Name  string  `json:"name"`
	Size  float64 `json:"size"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
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

type transfairChatMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Size    int    `json:"size"`
}

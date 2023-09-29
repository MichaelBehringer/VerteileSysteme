package main

import "github.com/google/uuid"

type Player struct {
	ID       uuid.UUID `json:"id"`
	Color    string    `json:"color"`
	Username string    `json:"name"`
}

type CustomPlayer struct {
	Skin     string `json:"skin"`
	Username string `json:"username"`
	Gamename string `json:"gamename"`
}

type NewPlayer struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Server struct {
	ID          uuid.UUID `json:"id"`
	PetName     string    `json:"petName"`
	Address     string    `json:"address"`
	PlayerCount int       `json:"playerCount"`
}

type AuthHeader struct {
	IDToken string `header:"Authorization"`
}

type Highscore struct {
	Highscore int    `json:"Highscore"`
	Name      string `json:"Name"`
}

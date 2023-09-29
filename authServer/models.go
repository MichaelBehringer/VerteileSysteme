package main

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Player struct {
	UserId string `json:"userId"`
}

type AccessToken struct {
	AccessToken string `json:"accessToken"`
}

type AuthHeader struct {
	IDToken string `header:"Authorization"`
}

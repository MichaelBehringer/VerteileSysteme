package main

import (
	"encoding/json"
	"flag"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/google/uuid"
)

// Loops over all Players and returns the top 5 with the highest scores
func getScore() []transfairScore {
	transfairScores := []transfairScore{}
	for _, value := range mapIdToPlayer {
		player := listPlayerKoordinates[value]
		transfairScores = append(transfairScores, transfairScore{Name: player.Name, Score: player.Size})
	}
	sort.Slice(transfairScores, func(i, j int) bool {
		return transfairScores[i].Score > transfairScores[j].Score
	})
	if len(transfairScores) >= 5 {
		return transfairScores[:5]
	} else {
		return transfairScores
	}
}

// get userid from token, and uservalues from database
func getTokenData(token string) (string, string, uuid.UUID) {
	// userId from token; Auth-Service needed for getting Value of Token
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://auth:8082/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, _ := client.Do(req)
	if res.StatusCode == 200 {
		body, _ := io.ReadAll(res.Body)
		var playerData map[string]interface{}
		json.Unmarshal(body, &playerData)
		userId, _ := playerData["userId"].(string)
		playerUUid, _ := uuid.Parse(userId)
		var player Player
		ExecuteSQLRow("SELECT Gamename, Skin FROM Player WHERE ID=?", userId).Scan(&player.Gamename, &player.Skin)
		return player.Gamename, player.Skin, playerUUid
	} else {
		return getRandomTokenData()
	}
}

// If user has no token, generate random username and color
func getRandomTokenData() (string, string, uuid.UUID) {
	username := getPetNameSingle()
	color := colors[rand.Intn(30)]
	return username, color, uuid.New()
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// get random server name
func getPetName() string {
	flag.Parse()
	return petname.Generate(2, "-")
}

// get random user name
func getPetNameSingle() string {
	return petname.Generate(1, "-")
}

// calls stored procedure for updating highscore; Logic when Highscore is updated is in the Database
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

// generating random coordinates and colors for initial NPCs
func initNPCs() {
	for i := 0; i < 1200; i++ {
		npc := gameObj{X: randFloat(0, mapBoundary), Y: randFloat(0, mapBoundary), Color: colors[rand.Intn(30)]}
		listNpcKoordinates[i] = npc
		treeNpc.Insert(Circle{Id: 0, X: npc.X, Y: npc.Y, Radius: 10, Color: npc.Color})
	}
}

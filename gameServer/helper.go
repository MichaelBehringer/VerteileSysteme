package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/dhconnelly/rtreego"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

type AccessToken struct {
	AccessToken string `json:"accessToken"`
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

func (c Circle) Bounds() *rtreego.Rect {
	return rtreego.Point{c.X, c.Y}.ToRect(c.Radius)
}

func npcCollision() {
	for _, value := range mapIdToPlayer {
		player := listPlayerKoordinates[value]
		searchCircle := Circle{Id: value, X: player.X, Y: player.Y, Radius: float64(player.Size)}

		results := treeNpc.SearchIntersect(searchCircle.Bounds())
		for _, result := range results {
			otherCircle := result.(Circle)
			if searchCircle != otherCircle {
				distance := math.Sqrt(math.Pow(searchCircle.X-otherCircle.X, 2) + math.Pow(searchCircle.Y-otherCircle.Y, 2))
				if distance < searchCircle.Radius+otherCircle.Radius {
					newX := randFloat(0, mapBoundary)
					newY := randFloat(0, mapBoundary)
					newColor := colors[rand.Intn(30)]
					listNpcKoordinates[otherCircle.Id].X = newX
					listNpcKoordinates[otherCircle.Id].Y = newY
					listNpcKoordinates[otherCircle.Id].Color = newColor
					currSize := listPlayerKoordinates[searchCircle.Id].Size
					if currSize > 500 {
						listPlayerKoordinates[searchCircle.Id].Size = currSize + (10.0 / currSize)
					} else {
						listPlayerKoordinates[searchCircle.Id].Size = currSize + (50.0 / currSize)
					}

					treeNpc.Delete(Circle{X: otherCircle.X, Y: otherCircle.Y, Radius: 10, Color: otherCircle.Color})
					treeNpc.Insert(Circle{X: newX, Y: newY, Radius: 10, Color: newColor})
				}
			}
		}
	}
}

func visibleNPC(playerObj transfairPlayer) []transfairNpc {
	var objNpcT []transfairNpc

	searchCircle := Circle{X: playerObj.X, Y: playerObj.Y, Radius: 500.0}

	results := treeNpc.SearchIntersect(searchCircle.Bounds())
	for _, result := range results {
		otherCircle := result.(Circle)
		objNpcT = append(objNpcT, transfairNpc{Color: otherCircle.Color, X: otherCircle.X, Y: otherCircle.Y})
	}
	return objNpcT
}

func getScore(playerObj transfairPlayer) []transfairScore {
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

func playerCollision() {
	tree := rtreego.NewTree(2, 25, 50)

	for _, value := range mapIdToPlayer {
		player := listPlayerKoordinates[value]
		tree.Insert(Circle{Id: value, X: player.X, Y: player.Y, Radius: float64(player.Size)})
	}

	for _, value := range mapIdToPlayer {
		player := listPlayerKoordinates[value]
		searchCircle := Circle{Id: value, X: player.X, Y: player.Y, Radius: float64(player.Size)}

		results := tree.SearchIntersect(searchCircle.Bounds())
		for _, result := range results {
			otherCircle := result.(Circle)
			if searchCircle != otherCircle {
				distance := math.Sqrt(math.Pow(searchCircle.X-otherCircle.X, 2) + math.Pow(searchCircle.Y-otherCircle.Y, 2))
				if distance < searchCircle.Radius+otherCircle.Radius {
					if listPlayerKoordinates[searchCircle.Id].Size > listPlayerKoordinates[otherCircle.Id].Size {
						listPlayerKoordinates[otherCircle.Id].X = randFloat(0, mapBoundary)
						listPlayerKoordinates[otherCircle.Id].Y = randFloat(0, mapBoundary)
						deadSize := listPlayerKoordinates[otherCircle.Id].Size
						listPlayerKoordinates[otherCircle.Id].Size = 20
						listPlayerKoordinates[searchCircle.Id].Size += (deadSize / 10)
					}
				}
			}
		}
	}
}

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

type Stack struct {
	items []int
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Push(item int) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() (int, error) {
	if s.IsEmpty() {
		return 0, fmt.Errorf("Stack is empty")
	}
	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]
	return item, nil
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func getTokenData(token string) (string, string, uuid.UUID) {
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

func getRandomTokenData() (string, string, uuid.UUID) {
	username := getPetNameSingle()
	color := colors[rand.Intn(30)]
	return username, color, uuid.New()
}

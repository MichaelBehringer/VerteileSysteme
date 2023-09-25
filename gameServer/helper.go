package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/dhconnelly/rtreego"
	"github.com/docker/docker/client"
)

type Circle struct {
	Id, Color int
	X, Y      float64
	Radius    float64
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
					newColor := rand.Intn(9)
					listNpcKoordinates[otherCircle.Id].X = newX
					listNpcKoordinates[otherCircle.Id].Y = newY
					listNpcKoordinates[otherCircle.Id].Color = newColor
					currSize := listPlayerKoordinates[searchCircle.Id].Size
					listPlayerKoordinates[searchCircle.Id].Size = currSize + (50.0 / currSize)

					treeNpc.Delete(Circle{X: otherCircle.X, Y: otherCircle.Y, Radius: 10, Color: otherCircle.Color})
					treeNpc.Insert(Circle{X: newX, Y: newY, Radius: 10, Color: newColor})
				}
			}
		}
	}
}

func visibleNPC(playerObj transfairPlayer) []transfaiNpc {
	var objNpcT []transfaiNpc

	searchCircle := Circle{X: playerObj.X, Y: playerObj.Y, Radius: 600.0}

	results := treeNpc.SearchIntersect(searchCircle.Bounds())
	for _, result := range results {
		otherCircle := result.(Circle)
		objNpcT = append(objNpcT, transfaiNpc{Color: otherCircle.Color, X: otherCircle.X, Y: otherCircle.Y})
	}
	return objNpcT
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
	containerName = strings.Split(containerName, "-")[2]
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

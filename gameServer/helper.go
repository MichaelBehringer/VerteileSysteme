package main

import (
	"fmt"
	"math"

	"github.com/dhconnelly/rtreego"
)

type Circle struct {
	Id     int
	X, Y   float64
	Radius float64
}

func (c Circle) Bounds() *rtreego.Rect {
	return rtreego.Point{c.X, c.Y}.ToRect(c.Radius)
}

func npcCollision() {
	tree := rtreego.NewTree(2, 25, 50)

	for id, value := range listNpcKoordinates {
		tree.Insert(Circle{Id: id, X: value.X, Y: value.Y, Radius: 10})
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
					fmt.Printf("NPC --- Kreise überlappen: (%f, %f), Radius %f und (%f, %f), Radius %f\n",
						searchCircle.X, searchCircle.Y, searchCircle.Radius,
						otherCircle.X, otherCircle.Y, otherCircle.Radius)
					listNpcKoordinates[otherCircle.Id].X = randFloat(0, 1000)
					listNpcKoordinates[otherCircle.Id].Y = randFloat(0, 700)
					listPlayerKoordinates[searchCircle.Id].Size += 2
				}
			}
		}
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
					fmt.Printf("Player --- Kreise überlappen: (%f, %f), Radius %f und (%f, %f), Radius %f\n",
						searchCircle.X, searchCircle.Y, searchCircle.Radius,
						otherCircle.X, otherCircle.Y, otherCircle.Radius)
					if listPlayerKoordinates[searchCircle.Id].Size > listPlayerKoordinates[otherCircle.Id].Size {
						listPlayerKoordinates[otherCircle.Id].X = randFloat(0, 1000)
						listPlayerKoordinates[otherCircle.Id].Y = randFloat(0, 700)
						listPlayerKoordinates[otherCircle.Id].Size = 20
						listPlayerKoordinates[searchCircle.Id].Size += 20
					}

				}
			}
		}
	}
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

package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/dhconnelly/rtreego"
)

// Calculate where Player moves, using 2d Vector and normalize it
func calcNewPoint(xStart float64, yStart float64, xEnd float64, yEnd float64, size float64) (float64, float64) {
	//map eingrenzen auf 0 bis 5k
	vectorX := xEnd - xStart
	vectorY := yEnd - yStart

	// Satz des Pythagoras
	lenghtVector := math.Sqrt(math.Pow(vectorX, 2) + math.Pow(vectorY, 2))
	if lenghtVector == 0 {
		return xStart, yStart
	}

	normalizedX := vectorX / lenghtVector
	normalizedY := vectorY / lenghtVector

	// Make Player move faster if he is small
	stepSize := 2.0 + 5.0/size

	// Move Player slower if the cursor is close to the middle
	if lenghtVector < 100 {
		stepSize = stepSize * lenghtVector * 0.01
	}

	newX := xStart + normalizedX*stepSize
	newY := yStart + normalizedY*stepSize

	// Prevent Player from leaving the map
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

func (c Circle) Bounds() *rtreego.Rect {
	return rtreego.Point{c.X, c.Y}.ToRect(c.Radius)
}

// Check Collision between Player and NPC using rTrees
func npcCollision() {
	for _, value := range mapIdToPlayer {
		player := listPlayerKoordinates[value]
		searchCircle := Circle{Id: value, X: player.X, Y: player.Y, Radius: float64(player.Size)}

		// Search for Collisions
		results := treeNpc.SearchIntersect(searchCircle.Bounds())
		for _, result := range results {
			otherCircle := result.(Circle)
			if searchCircle != otherCircle {
				// Check if Circles match; rtrees uses bounding boxes
				distance := math.Sqrt(math.Pow(searchCircle.X-otherCircle.X, 2) + math.Pow(searchCircle.Y-otherCircle.Y, 2))
				if distance < searchCircle.Radius+otherCircle.Radius {
					// if Collision, delete NPC and create new one; increase Player Size
					newX := randFloat(0, mapBoundary)
					newY := randFloat(0, mapBoundary)
					newColor := colors[rand.Intn(30)]
					listNpcKoordinates[otherCircle.Id].X = newX
					listNpcKoordinates[otherCircle.Id].Y = newY
					listNpcKoordinates[otherCircle.Id].Color = newColor
					currSize := listPlayerKoordinates[searchCircle.Id].Size
					// Slower growth if Player is big
					if currSize > 200 {
						listPlayerKoordinates[searchCircle.Id].Size = currSize + (10.0 / currSize)
					} else {
						listPlayerKoordinates[searchCircle.Id].Size = currSize + (50.0 / currSize)
					}

					// update npcTree to new Npc position
					treeNpc.Delete(Circle{X: otherCircle.X, Y: otherCircle.Y, Radius: 10, Color: otherCircle.Color})
					treeNpc.Insert(Circle{X: newX, Y: newY, Radius: 10, Color: newColor})
				}
			}
		}
	}
}

// Get all NPCs in a 500px radius; using rtrees agin for range
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

// Check Collision between Player and Player using rTrees
func playerCollision() {
	tree := rtreego.NewTree(2, 25, 50)

	// insert all Players into rTree
	for _, value := range mapIdToPlayer {
		player := listPlayerKoordinates[value]
		tree.Insert(Circle{Id: value, X: player.X, Y: player.Y, Radius: float64(player.Size)})
	}

	for _, value := range mapIdToPlayer {
		player := listPlayerKoordinates[value]
		searchCircle := Circle{Id: value, X: player.X, Y: player.Y, Radius: float64(player.Size)}

		// Search for Collisions
		results := tree.SearchIntersect(searchCircle.Bounds())
		for _, result := range results {
			otherCircle := result.(Circle)
			if searchCircle != otherCircle {
				// Check if Circles match; rtrees uses bounding boxes
				distance := math.Sqrt(math.Pow(searchCircle.X-otherCircle.X, 2) + math.Pow(searchCircle.Y-otherCircle.Y, 2))
				if distance < searchCircle.Radius+otherCircle.Radius {
					// if Collision, kill Player with smaller size
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

// Move Player every 10ms
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

// Check Collisions every 100ms
func checkCollission() {
	for range time.Tick(time.Second / 10) {
		npcCollision()
		playerCollision()
	}
}

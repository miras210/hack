package a_star_pathfinding_amiran

import (
	"fmt"
	"hackathon/models"
	"math"
	"strings"
)

const mapSize = 10000

func NewWorld(children []models.Coords, snowAreas []models.SnowArea) string {
	fmt.Println("Creating new world...")
	var gameMap [mapSize][mapSize]int

	for y := 0; y < mapSize; y++ {
		for x := 0; x < mapSize; x++ {
			curPoint := mapPoint{
				x: x,
				y: y,
			}
			for k := 0; k < len(snowAreas); k++ {
				ar := snowAreas[k]
				arPoint := mapPoint{
					x: ar.X,
					y: ar.Y,
				}
				if insideCircle(arPoint, curPoint, float64(ar.R)) {
					gameMap[y][x] = KindSnow
				}
			}
		}
	}

	var sb strings.Builder

	for y := 0; y < mapSize; y++ {
		for x := 0; x < mapSize; x++ {
			sb.WriteRune(KindRunes[gameMap[y][x]])
		}
		sb.WriteString("\n")
	}
	fmt.Println("World creation completed...")
	return sb.String()
}

func insideCircle(center, tile mapPoint, radius float64) bool {
	var dx = float64(center.x - tile.x)
	var dy = float64(center.y - tile.y)
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance <= radius
}

type mapPoint struct {
	x int
	y int
}

package solver

import "hackathon/models"

type Solver interface {
	Algo([]models.Coords, []models.Gift, []models.SnowArea) models.Request
}

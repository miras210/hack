package solver

import "hackathon/models"

type Solver interface {
	Algo([]models.Coords, []models.Gift, []models.SnowArea) models.Request
}

const MapID = "faf7ef78-41b3-4a36-8423-688a61929c08"

package solver

import "hackathon/models/phase2"

type Solver interface {
	Solve(data models.MapResponse) []models.PresentingGift
}

type DummySolver struct{}

func (d *DummySolver) Solve(data models.MapResponse) []models.PresentingGift {
	return nil
}

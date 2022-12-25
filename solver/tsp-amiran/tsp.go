package tsp_amiran

import (
	"hackathon/models"
	"sort"
)

var zeroPoint = models.Coords{
	X: 0,
	Y: 0,
}

func (g *TSPAmiranSolver) CalculateTSPMoves(children []models.Coords, bags int,
	x, y int) (newChildren, moves []models.Coords) {
	closestIdx := g.Closest(children, x, y)
	firstChild := children[closestIdx]
	children = removeFromSlice(children, closestIdx)
	if bags < 2 {
		newChildren = children
		moves = []models.Coords{
			{X: firstChild.X, Y: firstChild.Y},
			zeroPoint,
		}
		return
	}

	var subset []models.Coords

	children, subset = g.GetClosestSubsetToPoint(children, bags-1, firstChild)

	subset = append(subset, firstChild)

	var result []models.Coords
	var minPath float64
	var hasFirstResult bool
	hasFirstResult = false

	for i := 0; i < len(subset); i++ {
		target := subset[i]
		result, minPath, hasFirstResult = g.recursivePathCalculation(result, minPath, hasFirstResult, bags,
			g.Distance(target.X, target.Y, zeroPoint.X, zeroPoint.Y), []models.Coords{target}, subset, []int{i})
	}

	newChildren = children
	moves = result
	return
}

func (g *TSPAmiranSolver) recursivePathCalculation(result []models.Coords, minPath float64, hasFirstResult bool,
	maxDegree int, currentDistanceCost float64, currentCoordsPath []models.Coords, subset []models.Coords,
	passedIndices []int) ([]models.Coords, float64, bool) {

	if len(currentCoordsPath) == maxDegree {
		prev := subset[passedIndices[len(passedIndices)-1]]
		currentDistanceCost += g.Distance(prev.X, prev.Y, zeroPoint.X, zeroPoint.Y)
		currentCoordsPath = append(currentCoordsPath, models.Coords{
			X: 0,
			Y: 0,
		})

		if hasFirstResult {
			if minPath > currentDistanceCost {
				minPath = currentDistanceCost
				result = currentCoordsPath
			}
		} else {
			minPath = currentDistanceCost
			result = currentCoordsPath
			hasFirstResult = true
		}
		return result, minPath, hasFirstResult
	}

	for m := 0; m < len(subset); m++ {
		if contains(passedIndices, m) {
			continue
		}
		target := subset[m]
		prev := subset[passedIndices[len(passedIndices)-1]]
		result, minPath, hasFirstResult = g.recursivePathCalculation(result, minPath, hasFirstResult, maxDegree,
			currentDistanceCost+g.Distance(target.X, target.Y, prev.X, prev.Y), append(currentCoordsPath, target),
			subset, append(passedIndices, m))
	}
	return result, minPath, hasFirstResult
}

func contains(elems []int, v int) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func (g *TSPAmiranSolver) GetClosestSubsetToPoint(children []models.Coords, amountToExtract int,
	point models.Coords) (childrenFiltered []models.Coords,
	subset []models.Coords) {
	var coordPairs []coordPair

	for i := 0; i < len(children); i++ {
		ch := children[i]
		distance := g.Distance(point.X, point.Y, ch.X, ch.Y)
		coordPairs = append(coordPairs, coordPair{
			idx:      i,
			coord:    ch,
			distance: distance,
		})
	}

	sort.SliceStable(coordPairs, func(i, j int) bool {
		return coordPairs[i].distance < coordPairs[j].distance
	})

	coordPairs = coordPairs[:amountToExtract]

	var filteredCoords []models.Coords

	for i := 0; i < len(coordPairs); i++ {
		cp := coordPairs[i]
		children = removeFromSlice(children, cp.idx)
		for k := 0; k < len(coordPairs); k++ {
			if coordPairs[k].idx > cp.idx {
				coordPairs[k].idx--
			}

		}
		filteredCoords = append(filteredCoords, cp.coord)
	}

	childrenFiltered = children
	subset = filteredCoords

	return
}

type coordPair struct {
	idx      int
	coord    models.Coords
	distance float64
}

func removeFromSlice(slice []models.Coords, s int) []models.Coords {
	return append(slice[:s], slice[s+1:]...)
}

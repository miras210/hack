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
		remChildren := make([]models.Coords, len(subset))
		copy(remChildren, subset)
		target := remChildren[i]
		remChildren = removeFromSlice(remChildren, i)
		var path []models.Coords
		var dis = g.Distance(target.X, target.Y, zeroPoint.X, zeroPoint.Y)
		path = append(path, target)
		for m := 0; m < len(subset)-1; m++ {
			prev := target
			idx := g.Closest(remChildren, target.X, target.Y)
			target = remChildren[idx]
			path = append(path, target)
			dis += g.Distance(target.X, target.Y, prev.X, prev.Y)
			remChildren = removeFromSlice(remChildren, idx)
		}
		dis += g.Distance(target.X, target.Y, zeroPoint.X, zeroPoint.Y)
		path = append(path, models.Coords{
			X: 0,
			Y: 0,
		})
		if !hasFirstResult {
			result = path
			minPath = dis
		} else if minPath > dis {
			result = path
			minPath = dis
		}
	}

	newChildren = children
	moves = result
	return
}

func (g *TSPAmiranSolver) recursivePathCalculation(currentPath, remainingChildren []models.Coords,
	currentDistanceCost float64, prev models.Coords) ([]models.Coords, float64) {
	if len(remainingChildren) == 0 {
		currentDistanceCost += g.Distance(prev.X, prev.Y, zeroPoint.X, zeroPoint.Y)
		return currentPath, currentDistanceCost
	}

	idx := g.Closest(remainingChildren, prev.X, prev.Y)

	target := remainingChildren[idx]

	currentPath = append(currentPath, target)

	currentDistanceCost += g.Distance(prev.X, prev.Y, target.X, target.Y)

	remChildren := make([]models.Coords, len(remainingChildren))
	remChildren = append(remChildren, remainingChildren...)

	remChildren = removeFromSlice(remChildren, idx)

	return g.recursivePathCalculation(currentPath, remChildren, currentDistanceCost, target)

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

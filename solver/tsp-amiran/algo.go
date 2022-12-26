package tsp_amiran

import (
	"hackathon/models"
	"hackathon/solver"
	"math"
	"sort"
)

type TSPAmiranSolver struct {
	SnowAreas []models.SnowArea
}

type Bag struct {
	Gifts  []models.Gift
	Weight int
	Volume int
}

func (b *Bag) Result() []int {
	res := make([]int, 0, len(b.Gifts))
	for i := 0; i < len(b.Gifts); i++ {
		res = append(res, b.Gifts[len(b.Gifts)-i-1].ID)
	}
	return res
}
func (b *Bag) Add(gift models.Gift) bool {
	if b.Weight+gift.Weight > 200 {
		return false
	}
	if b.Volume+gift.Volume > 100 {
		return false
	}
	b.Gifts = append(b.Gifts, gift)
	b.Weight += gift.Weight
	b.Volume += gift.Volume
	return true
}
func (b *Bag) AddMax(gifts []models.Gift) int {
	for i := 0; i < len(gifts); i++ {
		if !b.Add(gifts[i]) {
			return i
		}
	}
	return len(gifts)
}

func (g *TSPAmiranSolver) Algo(children []models.Coords, gifts []models.Gift, snowAreas []models.SnowArea) models.Request {
	res := models.Request{
		MapID:       solver.MapID,
		Moves:       make([]models.Coords, 0, len(children)),
		StackOfBags: make([][]int, 0),
	}

	g.SnowAreas = snowAreas

	sort.SliceStable(gifts, func(i, j int) bool {
		return gifts[i].Weight+gifts[i].Volume <= gifts[j].Weight+gifts[j].Volume
	})

	for i, j := 0, len(gifts)-1; i < j; i, j = i+1, j-1 {
		gifts[i], gifts[j] = gifts[j], gifts[i]
	}

	for len(gifts) > 0 {
		currx, curry := 0, 0
		bag := Bag{
			Gifts:  make([]models.Gift, 0),
			Weight: 0,
			Volume: 0,
		}
		i := bag.AddMax(gifts)
		gifts = gifts[i:]
		res.StackOfBags = append(res.StackOfBags, bag.Result())

		count := len(bag.Gifts)
		//for count != 0 {
		//	idx := g.Closest(children, currx, curry)
		//	currx, curry = children[idx].X, children[idx].Y
		//	res.Moves = append(res.Moves, children[idx])
		//	children[idx] = children[len(children)-1]
		//	children = children[:len(children)-1]
		//	count--
		//}
		//if len(gifts) != 0 {
		//	zero := models.Coords{
		//		X: 0,
		//		Y: 0,
		//	}
		//	res.Moves = append(res.Moves, zero)
		//}

		var moves []models.Coords

		children, moves = g.CalculateTSPMoves(children, count, currx, curry)

		res.Moves = append(res.Moves, moves...)
	}
	return res
}
func (g *TSPAmiranSolver) Closest(children []models.Coords, x, y int) int {
	cx, cy := children[0].X, children[0].Y
	dist := g.Distance(cx, cy, x, y)
	ans := 0
	for i := 1; i < len(children); i++ {
		cx, cy = children[i].X, children[i].Y
		newDist := g.Distance(cx, cy, x, y)
		if newDist < dist {
			dist = newDist
			ans = i
		}
	}
	return ans
}
func Distance(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2))
}

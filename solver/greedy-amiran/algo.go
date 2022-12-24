package greedy_amiran

import (
	"hackathon/models"
	"hackathon/solver"
	"math"
)

const distanceMultiplier = 7

type GreedyAmiranSolver struct {
	SnowAreas []models.SnowArea
}

type Bag struct {
	Gifts  []models.Gift
	Weight int
	Volume int
}

type point struct{ x, y float64 }

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

func (g *GreedyAmiranSolver) Algo(children []models.Coords, gifts []models.Gift, snowAreas []models.SnowArea) models.Request {
	res := models.Request{
		MapID:       solver.MapID,
		Moves:       make([]models.Coords, 0, len(children)),
		StackOfBags: make([][]int, 0),
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
		for count != 0 {
			idx := g.Closest(children, currx, curry)
			currx, curry = children[idx].X, children[idx].Y
			res.Moves = append(res.Moves, children[idx])
			children[idx] = children[len(children)-1]
			children = children[:len(children)-1]
			count--
		}
		if len(gifts) != 0 {
			zero := models.Coords{
				X: 0,
				Y: 0,
			}
			res.Moves = append(res.Moves, zero)
		}
	}
	return res
}
func (g *GreedyAmiranSolver) Closest(children []models.Coords, x, y int) int {
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

func sq(x float64) float64 { return x * x }

type specialPoint struct {
	point
	isSnowy bool
}

func (g *GreedyAmiranSolver) Distance(x1, y1, x2, y2 int) float64 {
	startPoint := point{
		x: float64(x1),
		y: float64(y1),
	}
	endPoint := point{
		x: float64(x2),
		y: float64(y2),
	}

	var snowyDistance float64
	snowyDistance = 0

	for k := 0; k < len(g.SnowAreas); k++ {
		ar := g.SnowAreas[k]
		cp := point{
			x: float64(ar.X),
			y: float64(ar.Y),
		}

		var isStartSnowy = isInside(ar.X, ar.Y, ar.R, x1, y1)

		interPoints := intersects(startPoint, endPoint, cp, float64(ar.R), true)

		if len(interPoints) == 0 {
			continue
		} else if len(interPoints) == 1 {
			if isStartSnowy {
				snowyDistance += twoPointDistance(startPoint.x, startPoint.y, interPoints[0].x, interPoints[0].y)
			} else {
				snowyDistance += twoPointDistance(interPoints[0].x, interPoints[0].y, endPoint.x, endPoint.y)
			}
		} else if len(interPoints) == 2 {
			snowyDistance += twoPointDistance(interPoints[0].x, interPoints[0].y, interPoints[1].x, interPoints[1].y)
		}
	}

	normalDistance := twoPointDistance(startPoint.x, startPoint.y, endPoint.x, endPoint.y)

	normalDistance -= snowyDistance

	return normalDistance + (snowyDistance * distanceMultiplier)
}

func twoPointDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

const eps = 1e-14 // say

func intersects(p1, p2, cp point, r float64, segment bool) []point {
	var res []point
	x0, y0 := cp.x, cp.y
	x1, y1 := p1.x, p1.y
	x2, y2 := p2.x, p2.y
	A := y2 - y1
	B := x1 - x2
	C := x2*y1 - x1*y2
	a := sq(A) + sq(B)
	var b, c float64
	var bnz = true
	if math.Abs(B) >= eps { // if B isn't zero or close to it
		b = 2 * (A*C + A*B*y0 - sq(B)*x0)
		c = sq(C) + 2*B*C*y0 - sq(B)*(sq(r)-sq(x0)-sq(y0))
	} else {
		b = 2 * (B*C + A*B*x0 - sq(A)*y0)
		c = sq(C) + 2*A*C*x0 - sq(A)*(sq(r)-sq(x0)-sq(y0))
		bnz = false
	}
	d := sq(b) - 4*a*c // discriminant
	if d < 0 {
		// line & circle don't intersect
		return res
	}

	// checks whether a point is within a segment
	within := func(x, y float64) bool {
		d1 := math.Sqrt(sq(x2-x1) + sq(y2-y1)) // distance between end-points
		d2 := math.Sqrt(sq(x-x1) + sq(y-y1))   // distance from point to one end
		d3 := math.Sqrt(sq(x2-x) + sq(y2-y))   // distance from point to other end
		delta := d1 - d2 - d3
		return math.Abs(delta) < eps // true if delta is less than a small tolerance
	}

	var x, y float64
	fx := func() float64 { return -(A*x + C) / B }
	fy := func() float64 { return -(B*y + C) / A }
	rxy := func() {
		if !segment || within(x, y) {
			res = append(res, point{x, y})
		}
	}

	if d == 0 {
		// line is tangent to circle, so just one intersect at most
		if bnz {
			x = -b / (2 * a)
			y = fx()
			rxy()
		} else {
			y = -b / (2 * a)
			x = fy()
			rxy()
		}
	} else {
		// two intersects at most
		d = math.Sqrt(d)
		if bnz {
			x = (-b + d) / (2 * a)
			y = fx()
			rxy()
			x = (-b - d) / (2 * a)
			y = fx()
			rxy()
		} else {
			y = (-b + d) / (2 * a)
			x = fy()
			rxy()
			y = (-b - d) / (2 * a)
			x = fy()
			rxy()
		}
	}
	return res
}

func isInside(circleX, circle_y, rad, x, y int) bool {
	// Compare radius of circle with distance
	// of its center from given point
	if (x-circleX)*(x-circleX)+
		(y-circle_y)*(y-circle_y) <= rad*rad {
		return true
	} else {
		return false
	}
}

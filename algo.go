package main

import "math"

type Bag struct {
	Gifts  []Gift
	Weight int
	Volume int
}

func (b *Bag) Result() []int {
	res := make([]int, len(b.Gifts))

	for i := 0; i < len(b.Gifts); i++ {
		res[i] = b.Gifts[i].ID
	}

	return res
}

func (b *Bag) Add(gift Gift) bool {
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

func (b *Bag) AddMax(gifts []Gift) int {
	for i := len(gifts) - 1; i >= 0; i-- {
		if !b.Add(gifts[i]) {
			return i + 1
		}
	}

	return 0
}

func Algo(children []Coords, gifts []Gift) Request {
	res := Request{
		MapID:       "faf7ef78-41b3-4a36-8423-688a61929c08",
		Moves:       make([]Coords, 0, len(children)),
		StackOfBags: make([][]int, 0),
	}

	for len(gifts) != 0 {
		currx, curry := 0, 0
		bag := Bag{
			Gifts:  make([]Gift, 0),
			Weight: 0,
			Volume: 0,
		}

		idx := bag.AddMax(gifts)
		gifts = gifts[:idx]
		res.StackOfBags = append(res.StackOfBags, bag.Result())

		count := len(bag.Gifts)

		for count != 0 {
			idx := Closest(children, currx, curry)
			res.Moves = append(res.Moves, children[idx])
			children[idx] = children[len(children)-1]
			children = children[:len(children)-1]
			count--
		}

		if len(gifts) != 0 {
			zero := Coords{
				X: 0,
				Y: 0,
			}
			res.Moves = append(res.Moves, zero)
		}
	}

	return res
}

func Closest(children []Coords, x, y int) int {
	cx, cy := children[0].X, children[0].Y
	dist := Distance(cx, cy, x, y)
	ans := 0

	for i := 1; i < len(children); i++ {
		cx, cy = children[i].X, children[i].Y
		newDist := Distance(cx, cy, x, y)
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

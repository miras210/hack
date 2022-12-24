package a_star_pathfinding_amiran

import (
	"fmt"
	"github.com/beefsack/go-astar"
	"hackathon/models"
	"hackathon/solver"
	"log"
	"math"
	"strings"
)

type AStarPathfindingSolver struct {
	world     World
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

func (s *AStarPathfindingSolver) Algo(children []models.Coords, gifts []models.Gift,
	snowAreas []models.SnowArea) models.Request {
	res := models.Request{
		MapID:       solver.MapID,
		Moves:       make([]models.Coords, 0, len(children)),
		StackOfBags: make([][]int, 0),
	}

	s.world = ParseWorld(NewWorld(children, snowAreas))
	s.SnowAreas = snowAreas
	fmt.Println(res)

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
			idx, moves := s.Closest(children, currx, curry)
			currx, curry = children[idx].X, children[idx].Y
			res.Moves = append(res.Moves, moves...)
			children[idx] = children[len(children)-1]
			children = children[:len(children)-1]
			count--
		}
		if len(gifts) != 0 {
			zero := models.Coords{
				X: 0,
				Y: 0,
			}
			//res.Moves = append(res.Moves, zero)
			path, _, found := astar.Path(
				s.world.Tile(currx, curry),
				s.world.Tile(zero.X, zero.Y))
			if !found {
				res.Moves = append(res.Moves, zero)
			} else {
				res.Moves = append(res.Moves, pathToCoords(path)...)
			}
		}
	}
	return res
}

func (s *AStarPathfindingSolver) Closest(children []models.Coords, x, y int) (int, []models.Coords) {
	ans := 0
	for len(children) > 0 {
		var dist float64 = 0
		checked := false
		for i := 0; i < len(children); i++ {
			cx, cy := children[i].X, children[i].Y
			newDist := s.Distance(cx, cy, x, y)
			if !checked {
				checked = true
				dist = newDist
				ans = i
				continue
			}
			if newDist < dist {
				dist = newDist
				ans = i
			}
		}
		winnerChild := children[ans]
		path, distance, found := astar.Path(
			s.world.Tile(x, y),
			s.world.Tile(winnerChild.X, winnerChild.Y))
		if !found {
			log.Fatal("Something went terribly wrong")
		}
		if distance <= dist {
			return ans, pathToCoords(path)
		} else {
			children[ans] = children[len(children)-1]
			children = children[:len(children)-1]
			ans = 0
		}
	}
	return ans, []models.Coords{}
}

func pathToCoords(pathers []astar.Pather) []models.Coords {
	var coords []models.Coords
	for k := 0; k < len(pathers); k++ {
		tile := pathers[k].(*Tile)
		coords = append(coords, models.Coords{
			X: tile.X,
			Y: tile.Y,
		})
	}
	dirX, dirY := 0, 0
	for k := 0; k < len(coords)-1; {
		xd, yd := gridDistance(coords[k], coords[k+1])
		if xd == dirX && yd == dirY {
			coords = removeFromSlice(coords, k)
		} else {
			dirX = xd
			dirY = yd
			k++
		}
	}
	return reverseSlice(coords)[1:]
}

func dot(x1, y1, x2, y2 int) int {
	sum := 0
	sum += x1 * x2
	sum += y1 * y2

	return sum
}

func removeFromSlice(slice []models.Coords, s int) []models.Coords {
	return append(slice[:s], slice[s+1:]...)
}

func reverseSlice(s []models.Coords) []models.Coords {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func gridDistance(from, to models.Coords) (x, y int) {
	return to.X - from.X, to.Y - from.Y
}

func Distance(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2))
}

// Kind* constants refer to tile kinds for input and output.
const (
	// KindPlain (.) is a plain tile with a movement cost of 1.
	KindPlain = iota
	// KindSnow (*) is a river tile with a movement cost of 7.
	KindSnow
	// KindFrom (F) is a tile which marks where the path should be calculated
	// from.
	KindFrom
	// KindTo (T) is a tile which marks the goal of the path.
	KindTo
	// KindPath (●) is a tile to represent where the path is in the output.
	KindPath
)

// KindRunes map tile kinds to output runes.
var KindRunes = map[int]rune{
	KindPlain: '.',
	KindSnow:  '*',
	KindFrom:  'F',
	KindTo:    'T',
	KindPath:  '●',
}

// RuneKinds map input runes to tile kinds.
var RuneKinds = map[rune]int{
	'.': KindPlain,
	'❆': KindSnow,
	'F': KindFrom,
	'T': KindTo,
}

// KindCosts map tile kinds to movement costs.
var KindCosts = map[int]float64{
	KindPlain: 1.0,
	KindSnow:  7.0,
	KindFrom:  1.0,
	KindTo:    1.0,
}

// PathNeighbors returns the neighbors of the tile, excluding blockers and
// tiles off the edge of the board.
func (t *Tile) PathNeighbors() []astar.Pather {
	var neighbors []astar.Pather
	for _, offset := range [][]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
		{-1, -1},
		{1, 1},
		{1, -1},
		{-1, 1},
	} {
		if n := t.W.Tile(t.X+offset[0], t.Y+offset[1]); n != nil {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

// A Tile is a tile in a grid which implements Pather.
type Tile struct {
	// Kind is the kind of tile, potentially affecting movement.
	Kind int
	// X and Y are the coordinates of the tile.
	X, Y int
	// W is a reference to the World that the tile is a part of.
	W World
}

// PathNeighborCost returns the movement cost of the directly neighboring tile.
func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	toT := to.(*Tile)
	return KindCosts[toT.Kind]
}

// PathEstimatedCost uses Manhattan distance to estimate orthogonal distance
// between non-adjacent nodes.
func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(*Tile)
	absX := toT.X - t.X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Y - t.Y
	if absY < 0 {
		absY = -absY
	}
	return math.Sqrt(math.Pow(float64(absX), 2) + math.Pow(float64(absY), 2))
}

// World is a two dimensional map of Tiles.
type World map[int]map[int]*Tile

// Tile gets the tile at the given coordinates in the world.
func (w World) Tile(x, y int) *Tile {
	if w[x] == nil {
		return nil
	}
	return w[x][y]
}

// SetTile sets a tile at the given coordinates in the world.
func (w World) SetTile(t *Tile, x, y int) {
	if w[x] == nil {
		w[x] = map[int]*Tile{}
	}
	w[x][y] = t
	t.X = x
	t.Y = y
	t.W = w
}

// FirstOfKind gets the first tile on the board of a kind, used to get the from
// and to tiles as there should only be one of each.
func (w World) FirstOfKind(kind int) *Tile {
	for _, row := range w {
		for _, t := range row {
			if t.Kind == kind {
				return t
			}
		}
	}
	return nil
}

// From gets the from tile from the world.
func (w World) From() *Tile {
	return w.FirstOfKind(KindFrom)
}

// To gets the to tile from the world.
func (w World) To() *Tile {
	return w.FirstOfKind(KindTo)
}

// RenderPath renders a path on top of a world.
func (w World) RenderPath(path []astar.Pather) string {
	width := len(w)
	if width == 0 {
		return ""
	}
	height := len(w[0])
	pathLocs := map[string]bool{}
	for _, p := range path {
		pT := p.(*Tile)
		pathLocs[fmt.Sprintf("%d,%d", pT.X, pT.Y)] = true
	}
	rows := make([]string, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			t := w.Tile(x, y)
			r := ' '
			if pathLocs[fmt.Sprintf("%d,%d", x, y)] {
				r = KindRunes[KindPath]
			} else if t != nil {
				r = KindRunes[t.Kind]
			}
			rows[y] += string(r)
		}
	}
	return strings.Join(rows, "\n")
}

// ParseWorld parses a textual representation of a world into a world map.
func ParseWorld(input string) World {
	fmt.Println("Started world parsing...")
	w := World{}
	for y, row := range strings.Split(strings.TrimSpace(input), "\n") {
		for x, raw := range row {
			kind, _ := RuneKinds[raw]
			w.SetTile(&Tile{
				Kind: kind,
			}, x, y)
		}
	}
	fmt.Println("Completed world parsing...")
	return w
}

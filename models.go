package main

type Response struct {
	Gifts     []Gift     `json:"gifts"`
	SnowAreas []SnowArea `json:"snowAreas"`
	Children  []Coords   `json:"children"`
}

type Gift struct {
	ID     int `json:"id"`
	Weight int `json:"weight"`
	Volume int `json:"volume"`
}

type SnowArea struct {
	R int `json:"r"`
	X int `json:"x"`
	Z int `json:"z"`
}

type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Request struct {
	MapID       string   `json:"mapID"`
	Moves       []Coords `json:"moves"`
	StackOfBags [][]int  `json:"stackOfBags"`
}

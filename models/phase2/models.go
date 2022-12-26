package models

type MapResponse struct {
	Gifts    []Gift  `json:"gifts"`
	Children []Child `json:"children"`
}

type Child struct {
	ID     int    `json:"id"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

type Gift struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Price int    `json:"price"`
}

type PresentingGift struct {
	GiftID  int `json:"giftID"`
	ChildID int `json:"childID"`
}

type Request struct {
	MapID           string           `json:"mapID"`
	PresentingGifts []PresentingGift `json:"presentingGifts"`
}

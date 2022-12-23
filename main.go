package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("response.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	byteVal, _ := io.ReadAll(f)
	var input Response
	json.Unmarshal(byteVal, &input)

	fmt.Println(Algo(input.Children, input.Gifts))
}

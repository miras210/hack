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

	res := Algo(input.Children, input.Gifts)
	fmt.Println(len(input.Children))
	fmt.Println(len(res.Moves))
	fmt.Println(res.StackOfBags)
	/*b, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	bodyReader := bytes.NewReader(b)

	req, err := http.NewRequest(http.MethodPost, "https://datsanta.dats.team/api/round", bodyReader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-API-Key", "92810ac8-2890-4b01-9379-151be16fbbee")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Обработка ошибки
		return
	}
	defer resp.Body.Close()*/
}

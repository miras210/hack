package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	request := Algo(input.Children, input.Gifts)

	body, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	url := "https://datsanta.dats.team/api/round"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-API-Key", "92810ac8-2890-4b01-9379-151be16fbbee")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	//print body in strings
	fmt.Println(string(body))
}

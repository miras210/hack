package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	res := Algo(input.Children, input.Gifts)
	b, err := json.Marshal(res)
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
	defer resp.Body.Close()
	var response Resp
	json.NewDecoder(resp.Body).Decode(response)
	fmt.Println(response)

	moves, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	fe, err := os.Create("./visual/moves.json")

	if err != nil {
		log.Fatal(err)
	}

	defer fe.Close()

	_, err2 := fe.Write(moves)

	if err2 != nil {
		log.Fatal(err2)
	}

}

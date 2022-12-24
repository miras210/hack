package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"hackathon/models"
	solver2 "hackathon/solver"
	"hackathon/solver/greedy-miras"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	var flagTest bool

	flag.BoolVar(&flagTest, "test", false, "Include if only for testing")

	flag.Parse()

	f, err := os.Open("response.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	byteVal, _ := io.ReadAll(f)
	var input models.Response
	json.Unmarshal(byteVal, &input)

	var solverImpl solver2.Solver

	//CHANGEABLE BLOCK

	solverImpl =
		&greedy_miras.GreedyMirasSolver{}

	//

	request := solverImpl.Algo(input.Children, input.Gifts, input.SnowAreas)

	for i, j := 0, len(request.StackOfBags)-1; i < j; i, j = i+1, j-1 {
		request.StackOfBags[i], request.StackOfBags[j] = request.StackOfBags[j], request.StackOfBags[i]
	}

	fmt.Println(len(input.Children))
	fmt.Println(len(request.Moves))
	fmt.Println(request.StackOfBags)
	body, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	if !flagTest {
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

		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		//print body in strings
		fmt.Println(string(respbody))
	}

	fe, err := os.Create("./visual/moves.json")

	if err != nil {
		log.Fatal(err)
	}

	defer fe.Close()

	_, err2 := fe.Write(body)

	if err2 != nil {
		log.Fatal(err2)
	}

}

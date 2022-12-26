package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"hackathon/models/phase2"
	"hackathon/phase2/solver"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	var flagTest bool

	flag.BoolVar(&flagTest, "test", false, "Include if only for testing")

	flag.Parse()

	f, err := os.Open("response2.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	byteVal, _ := io.ReadAll(f)
	var input models.MapResponse
	json.Unmarshal(byteVal, &input)

	var solverImpl solver.Solver
	solverImpl =
		//CHANGE THIS CODE HERE
		&solver.DummySolver{}
	//

	solvedData := solverImpl.Solve(input)

	if solvedData == nil {
		log.Fatal("Solver returned nil response! Or maybe you are using DummySolver?")
	}

	fmt.Printf("Length of PresentingGifts: %d\n", len(solvedData))

	request := models.Request{
		MapID:           "a8e01288-28f8-45ee-9db4-f74fc4ff02c8",
		PresentingGifts: solvedData,
	}

	body, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	if !flagTest {
		url := "https://datsanta.dats.team/api/round2"
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
}

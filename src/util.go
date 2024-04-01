package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type BusRouters struct {
	Categories []struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Routes []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"routes"`
	} `json:"categories"`
}

func getBusRoutes() BusRouters {
	filepath := filepath.Join("src", "url.json")
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var bus BusRouters
	if err := json.NewDecoder(f).Decode(&bus); err != nil {
		log.Fatal(err)
	}

	return bus
}

func sortBusResponse(busResponse BusResponse) BusResponse {
	for _, buses := range busResponse.BusTimetables {
		for i := 0; i < len(buses); i++ {
			for j := i + 1; j < len(buses); j++ {
				if buses[i].OnTime > buses[j].OnTime {
					buses[i], buses[j] = buses[j], buses[i]
				}
			}
		}
	}

	return busResponse
}

package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/cors"
)

var (
	busTimetables     BusResponse
	busTimetablesTime time.Time
)

const cacheExpire = 60 * time.Second

func server() {
	http.HandleFunc("GET /v1/all", serverHandler)

	c := cors.Default()
	handler := c.Handler(http.DefaultServeMux)

	http.ListenAndServe(":8080", handler)
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	if time.Since(busTimetablesTime) < cacheExpire {
		busTimetables.IsCached = true
	} else {
		busTimetables = getBusTimetables()
		busTimetablesTime = time.Now()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(busTimetables)
}

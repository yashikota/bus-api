package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/rs/cors"
)

var (
	busTimetables     BusResponse
	busTimetablesTime time.Time
	busTimetablesLock sync.RWMutex
)

const cacheExpire = 70 * time.Second

func server() {
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "robots.txt")
	})
	http.HandleFunc("/v1/all", serverHandler)

	c := cors.Default()
	handler := c.Handler(http.DefaultServeMux)

	go refreshBusTimetables()

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

func refreshBusTimetables() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		newBusTimetables := getBusTimetables()
		busTimetablesLock.Lock()
		busTimetables = newBusTimetables
		busTimetablesTime = time.Now()
		busTimetablesLock.Unlock()
	}
}

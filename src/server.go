package main

import (
	"encoding/json"
	"net/http"
)

func server() {
	http.HandleFunc("GET /v1/all", serverHandler)
	http.ListenAndServe(":8080", nil)
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	busTimetables := getBusTimetables()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(busTimetables)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var busTimetables BusResponse
var busTimetablesCache BusResponse
var busTimetablesCacheDate int64

func server() {
	http.HandleFunc("GET /v1/all", serverHandler)
	http.ListenAndServe(":8080", nil)
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	// キャッシュが存在し、60秒以内ならキャッシュを返す
	if busTimetablesCacheDate > 0 && time.Now().Unix()-busTimetablesCacheDate < 60 {
		fmt.Println("Cache hit")
		busTimetablesCache.IsCached = true
		busTimetables = busTimetablesCache
	} else {
		fmt.Println("Cache miss")
		busTimetables = getBusTimetables()
		busTimetablesCache = busTimetables
		busTimetablesCacheDate = time.Now().Unix()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(busTimetables)
}

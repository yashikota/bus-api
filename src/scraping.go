package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

type BusResponse struct {
	IsCached      bool
	FetchTime     string
	BusTimetables map[string][]Bus
}

type Bus struct {
	BusStop       string
	Stand         string
	Name          string
	IsSignal      bool
	OnTime        string
	EstimatedTime string
	MoreMinutes   string
	DelayMinutes  int
	System        string
	Destination   string
}

func extract(s, pattern string) string {
	return regexp.MustCompile(pattern).FindString(s)
}

func extractMatch(s, pattern string) string {
	s = regexp.MustCompile(`\s`).ReplaceAllString(s, "")
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(s)

	return match[1]
}

func scrape(url string, busStop string, stand string, name string) []Bus {
	var buses []Bus
	c := colly.NewCollector()

	c.OnHTML(".pc.busstateArea", func(e *colly.HTMLElement) {
		e.ForEach(".pc .divbusstate", func(_ int, e *colly.HTMLElement) {
			// 非表示にされている経路をスキップ
			if e.Attr("class") == "pc divbusstate notview" {
				return
			}

			// 定刻
			fixedTime := extract(e.ChildText(".bsul.first"), `\b(\d{2}:\d{2})\b`)

			// 予定時間 or 定刻
			estimatedTime := extract(e.ChildText(".time"), `\b(\d{2}:\d{2})\b`)

			// 遅延時間
			var delayMinutes int
			capture := extract(e.ChildText(".bsul.first"), `\(.*\)`)
			// 遅延していた場合
			if capture != "(定時運行中)" {
				delay, _ := strconv.Atoi(extract(capture, `(\d+)`))
				delayMinutes = delay
			}

			// 残り時間
			moreMinutes := extract(e.ChildText(".more_min"), `(\d+時間)?(\d+分)|まもなく到着`)

			// 受信中 or 未受信
			isSignal := e.ChildText(".signal_status") != ""

			// 系統
			system := extractMatch(e.ChildText(".bsul"), `系統：\[(.*?)\]行先`)

			// 行先
			destination := extractMatch(e.ChildText(".bsul"), `行先：(.*?)行`)

			fmt.Println("BusStop:", busStop, "Stand:", stand, "Name:", name, "IsSignal:", isSignal, "OnTime:", fixedTime, "EstimatedTime:", estimatedTime, "MoreMinutes:", moreMinutes, "DelayMinutes:", delayMinutes, "System:", system, "Destination:", destination)

			bus := Bus{
				BusStop:       busStop,
				Stand:         stand,
				Name:          name,
				IsSignal:      isSignal,
				OnTime:        fixedTime,
				EstimatedTime: estimatedTime,
				MoreMinutes:   moreMinutes,
				DelayMinutes:  delayMinutes,
				System:        system,
				Destination:   destination,
			}
			buses = append(buses, bus)
		})
	})
	c.Visit(url)

	return buses
}

func getBusTimetables() BusResponse {
	busRouters := getBusRoutes()
	var keyNames = []string{"Kuzuha-OIT", "Nagao-OIT", "OIT-Kuzuha", "OIT-Nagao"}
	busResponse := BusResponse{
		BusTimetables: make(map[string][]Bus),
		FetchTime:     time.Now().Format("15:04"), // mm:ss
	}

	for i, category := range busRouters.Categories {
		for _, route := range category.Routes {
			fmt.Println(category.From, "=>", category.To, ":", route.Name)
			buses := scrape(route.URL, category.From, category.To, route.Name)
			busResponse.BusTimetables[keyNames[i]] = append(busResponse.BusTimetables[keyNames[i]], buses...)
		}

		fmt.Println("=========================================")
		time.Sleep(1 * time.Second)
	}

	return sortBusResponse(busResponse)
}

package main

import (
	"fmt"
	"net/http"
	"io"
	"os"
	"time"
	"encoding/json"
)

var coordinates = [][]string{
	{"43.766599", "-79.347110"},
	{"43.766927", "-79.346161"},
	{"43.768203", "-79.330267"},
	{"43.768574", "-79.32774"},
	{"43.751075", "-79.411487"},
	{"43.751206", "-79.411792"},
	{"43.756596", "-79.403697"},
	{"43.756871", "-79.40369"},
	{"43.728015", "-79.456212"},
	{"43.728259", "-79.456219"},
	{"43.734564", "-79.436478"},
	{"43.735002", "-79.43604"},
	{"43.714708", "-79.528720"},
	{"43.714819", "-79.529243"},
	{"43.717437", "-79.511101"},
	{"43.717569", "-79.51133"},
	{"43.665526", "-79.598759"},
	{"43.665706", "-79.599038"},
	{"43.676536", "-79.577802"},
	{"43.676426", "-79.578167"},
}

type ApiResponse struct {
	SpeedInformation SpeedInformation `json:"flowSegmentData"`
}

type SpeedInformation struct {
	CurrentSpeed float64 `json:"currentSpeed"`
	FreeFlowSpeed float64 `json:"freeFlowSpeed"`
}

func fetchAPIBody() {
	for _, latlng := range coordinates{
		resp, err := http.Get(fmt.Sprintf("https://api.tomtom.com/traffic/services/4/flowSegmentData/relative0/10/json?point=%s%%2C%s&unit=KMPH&openLr=false&key=%s", latlng[0], latlng[1], os.Getenv("TOMTOM_API_KEY")))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// Handle error reading body
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}
		// fmt.Printf("Response body: %s\n", body)

		var apiResponse ApiResponse
		err = json.Unmarshal(body, &apiResponse)
		if (err != nil) {
			fmt.Printf("Error unmarshalling JSON: %v\n", err)
			return
		}

		var congestionLevel string = ""

		if (apiResponse.SpeedInformation.CurrentSpeed <= 0.3 * apiResponse.SpeedInformation.FreeFlowSpeed) {
			congestionLevel = "High Congestion"
		} else if (apiResponse.SpeedInformation.CurrentSpeed <= 0.7 * apiResponse.SpeedInformation.FreeFlowSpeed) {
			congestionLevel = "Medium Congestion"
		} else {
			congestionLevel = "No Congestion"
		}

		fmt.Printf("For lat,lng (%s, %s): Current Speed: %f, Free Flow Speed: %f --> %s\n", latlng[0], latlng[1], apiResponse.SpeedInformation.CurrentSpeed, apiResponse.SpeedInformation.FreeFlowSpeed, congestionLevel)
	}


}

func main() {
	ticket := time.NewTicker(5 * time.Second)
	defer ticket.Stop()
	fetchAPIBody()

	for {
		select {
		case <-ticket.C:
			fmt.Printf("Ticker ticked at %v\n", time.Now())
			fetchAPIBody()
		}
	}
}

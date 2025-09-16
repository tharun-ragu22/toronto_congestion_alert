package main

import (
	"fmt"
	"net/http"
	"io"
	"os"
	"time"
	"encoding/json"
)

type ApiResponse struct {
	SpeedInformation SpeedInformation `json:"flowSegmentData"`
}

type SpeedInformation struct {
	CurrentSpeed float64 `json:"currentSpeed"`
	FreeFlowSpeed float64 `json:"freeFlowSpeed"`
}

func fetchAPIBody() {
	resp, err := http.Get(fmt.Sprintf("https://api.tomtom.com/traffic/services/4/flowSegmentData/relative0/10/json?point=43.729889%%2C-79.447357&unit=KMPH&openLr=false&key=%s", os.Getenv("TOMTOM_API_KEY")))
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

	fmt.Printf("Current Speed: %f, Free Flow Speed: %f\n", apiResponse.SpeedInformation.CurrentSpeed, apiResponse.SpeedInformation.FreeFlowSpeed)


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

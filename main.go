package main

import (
	"fmt"
	"net/http"
	"io"
	"os"
	"encoding/json"
	"log"
	"time"
)

type Exit struct {
	Name       string
	ExitPoints []GPS
}

type GPS struct {
	Lat  string
	Long string
}

var exits = []Exit{
	{
		Name: "DVP",
		ExitPoints: []GPS{
			{Lat: "43.766599", Long: "-79.347110"},
			{Lat: "43.766927", Long: "-79.346161"},
			{Lat: "43.768203", Long: "-79.330267"},
			{Lat: "43.768574", Long: "-79.327740"},
		},
	},
	{
		Name: "Yonge",
		ExitPoints: []GPS{
			{Lat: "43.751075", Long: "-79.411487"},
			{Lat: "43.751206", Long: "-79.411792"},
			{Lat: "43.756596", Long: "-79.403697"},
			{Lat: "43.756871", Long: "-79.403690"},
		},
	},
	{
		Name: "Allen",
		ExitPoints: []GPS{
			{Lat: "43.728015", Long: "-79.456212"},
			{Lat: "43.728259", Long: "-79.456219"},
			{Lat: "43.734564", Long: "-79.436478"},
			{Lat: "43.735002", Long: "-79.436040"},
		},
	},
	{
		Name: "Hwy 400",
		ExitPoints: []GPS{
			{Lat: "43.714708", Long: "-79.528720"},
			{Lat: "43.714819", Long: "-79.529243"},
			{Lat: "43.717437", Long: "-79.511101"},
			{Lat: "43.717569", Long: "-79.511330"},
		},
	},
	{
		Name: "Hwy 427",
		ExitPoints: []GPS{
			{Lat: "43.665526", Long: "-79.598759"},
			{Lat: "43.665706", Long: "-79.599038"},
			{Lat: "43.676536", Long: "-79.577802"},
			{Lat: "43.676426", Long: "-79.578167"},
		},
	},
}

var exitOrder = []string{
	"Eastbound Collectors",
	"Eastbound Express",
	"Westbound Collectors",
	"Westbound Express",
}

type ApiResponse struct {
	SpeedInformation SpeedInformation `json:"flowSegmentData"`
}

type SpeedInformation struct {
	CurrentSpeed float64 `json:"currentSpeed"`
	FreeFlowSpeed float64 `json:"freeFlowSpeed"`
}

func fetchAPIBody(w http.ResponseWriter) {
	for _, exit := range exits{
		for i, exitCoordinates := range exit.ExitPoints{
			resp, err := http.Get(fmt.Sprintf("https://api.tomtom.com/traffic/services/4/flowSegmentData/relative0/10/json?point=%s%%2C%s&unit=KMPH&openLr=false&key=%s", exitCoordinates.Lat, exitCoordinates.Long, os.Getenv("TOMTOM_API_KEY")))
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

			if (apiResponse.SpeedInformation.CurrentSpeed <= 0.5 * apiResponse.SpeedInformation.FreeFlowSpeed) {
				congestionLevel = "High Congestion"
			} else if (apiResponse.SpeedInformation.CurrentSpeed <= 0.9 * apiResponse.SpeedInformation.FreeFlowSpeed) {
				congestionLevel = "Medium Congestion"
			} else {
				congestionLevel = "No Congestion"
			}

			fmt.Fprintf(w,"401 and %s %s: Current Speed: %f, Free Flow Speed: %f --> %s\n", exit.Name, exitOrder[i], apiResponse.SpeedInformation.CurrentSpeed, apiResponse.SpeedInformation.FreeFlowSpeed, congestionLevel)
		}
	}


}

// func main() {
// 	ticket := time.NewTicker(5 * time.Second)
// 	defer ticket.Stop()
// 	fetchAPIBody()

// 	for t := range ticket.C {
// 		fmt.Printf("Ticker ticked at %v\n", t)
// 		fetchAPIBody()
// 	}
// }

func handler(w http.ResponseWriter, r *http.Request) {
		ticket := time.NewTicker(5 * time.Second)
		defer ticket.Stop()
		fetchAPIBody(w)

		for t := range ticket.C {
			fmt.Fprintf(w, "Ticker ticked at %v\n", t)
			fetchAPIBody(w)
		}
}

func main() {
	port := os.Getenv("PORT") // Render sets PORT automatically
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", handler)
	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

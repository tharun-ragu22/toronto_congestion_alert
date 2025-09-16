package main

import (
	"fmt"
	"net/http"
	"io"
	"os"
)

func main() {
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
    fmt.Printf("Response body: %s\n", body)
}

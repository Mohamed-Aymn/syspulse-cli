package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"syspulse-cli/identifiers"
	"syspulse-cli/metrics"
	"syspulse-cli/websocket"
)

func main() {
	key := flag.String("key", "", "provide the key that is generated from the website")
	flag.Parse()

	if *key == "" {
		fmt.Println("Error: -key flag is required")
		os.Exit(1)
	}

	isNew := !identifiers.IsIDStored()

	// Validate key and capture the response
	resp, err := identifiers.ValidateKey(*key, isNew)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close() // Ensure response body is closed

	// Fetch and store ID if needed
	if isNew {
		id, err := identifiers.FetchAndStoreID(resp)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("ID saved successfully:", id)
	}

	// Fetch stored device ID
	deviceID, err := identifiers.ReadIDFromFile()
	if err != nil {
		fmt.Println("Error retrieving device ID:", err)
		os.Exit(1)
	}

	// WebSocket connection
	wsURL := "ws://localhost:3000"
	client, err := websocket.NewWebSocketClient(wsURL)
	if err != nil {
		fmt.Printf("WebSocket connection failed: %v\n", err)
		return
	}
	defer client.Close()

	// Initialize lastMetrics as an empty map
	lastMetrics := make(map[string]string)
	sleepTime := 5 * time.Second
	skipCounter := 0 // the number of times metrics are calculated, skipped and not sent
	skipLimit := 12  // max number of times to skip sending metrics to the server (n * sleep time = total time)

	for {
		currentMetrics := metrics.GetMetrics(deviceID)

		if metrics.HasSignificantChange(lastMetrics, currentMetrics, 5.0) || skipCounter >= skipLimit {
			if err := client.SendMessage(currentMetrics); err != nil {
				fmt.Println("Error sending WebSocket message:", err)
			} else {
				lastMetrics = currentMetrics
				skipCounter = 0
			}
		} else {
			skipCounter++
			fmt.Println("No significant change in metrics.")
		}

		time.Sleep(sleepTime)
	}
}

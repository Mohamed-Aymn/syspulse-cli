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

	for {
		currentMetrics := metrics.GetMetrics()

		if metrics.HasSignificantChange(lastMetrics, currentMetrics, 5.0) {
			if err := client.SendMessage(currentMetrics); err != nil {
				fmt.Println("Error sending WebSocket message:", err)
			} else {
				lastMetrics = currentMetrics
			}
		} else {
			fmt.Println("No significant change in metrics.")
		}

		time.Sleep(sleepTime)
	}
}

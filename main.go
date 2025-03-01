package main

import (
	"fmt"
	"time"

	"syspulse-cli/metrics"
	"syspulse-cli/websocket"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	wsURL := "ws://localhost:8080"
	client, err := websocket.NewWebSocketClient(wsURL)
	if err != nil {
		fmt.Printf("WebSocket connection failed: %v\n", err)
		return
	}
	defer client.Close()

	var lastMetrics map[string]string
	sleepTime := 5 * time.Second
	maxSleepTime := 60 * time.Second

	for {
		currentMetrics := metrics.GetMetrics()

		if lastMetrics == nil || metrics.HasSignificantChange(lastMetrics, currentMetrics, 5.0) {
			if err := client.SendMessage(currentMetrics); err != nil {
				fmt.Println(err)
			} else {
				sleepTime = 5 * time.Second // Reset sleep time
				lastMetrics = currentMetrics
			}
		} else {
			sleepTime = time.Duration(min(int64(sleepTime)*2, int64(maxSleepTime)))
			fmt.Printf("Sleep time doubled to %v\n", sleepTime)
		}

		time.Sleep(sleepTime)
	}
}

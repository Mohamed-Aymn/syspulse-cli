package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func runCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return strings.TrimSpace(string(output)) // Trim spaces & newlines
}

func getMetrics() map[string]string {
	commands := map[string][]string{
		"CPU Usage":    {"sh", "-c", "grep 'cpu ' /proc/stat | awk '{usage=100-($5*100/($2+$3+$4+$5+$6+$7+$8))} END {print usage \"%\"}'"},
		"Memory Usage": {"sh", "-c", "free | awk '/Mem:/ {print ($3/$2)*100 \"%\"}'"},
	}

	metrics := make(map[string]string)
	for name, args := range commands {
		metrics[name] = runCommand(args[0], args[1:]...)
	}

	return metrics
}

func hasSignificantChange(old, new map[string]string, threshold float64) bool {
	for key, newValue := range new {
		oldValue, exists := old[key]
		if !exists {
			return true
		}

		// Remove "%" sign and convert to float
		oldFloat, err1 := strconv.ParseFloat(strings.TrimSuffix(oldValue, "%"), 64)
		newFloat, err2 := strconv.ParseFloat(strings.TrimSuffix(newValue, "%"), 64)

		if err1 != nil || err2 != nil {
			return true // If conversion fails, assume change is significant
		}

		if abs(newFloat-oldFloat) >= threshold {
			return true
		}
	}
	return false
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

func main() {
	wsURL := "ws://localhost:8080"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		fmt.Printf("WebSocket connection failed: %v\n", err)
		return
	}
	defer conn.Close()

	var lastMetrics map[string]string
	sleepTime := 5 * time.Second
	maxSleepTime := 60 * time.Second

	for {
		currentMetrics := getMetrics()

		if lastMetrics == nil || hasSignificantChange(lastMetrics, currentMetrics, 5.0) {
			jsonData, err := json.Marshal(currentMetrics)
			if err != nil {
				fmt.Printf("Error encoding JSON: %v\n", err)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				fmt.Printf("Error sending message: %v\n", err)
			} else {
				fmt.Printf("Sent: %s\n", jsonData)
			}

			// Reset sleep time if change is detected
			sleepTime = 5 * time.Second
			lastMetrics = currentMetrics
		} else {
			// Double sleep time up to max limit
			sleepTime = time.Duration(min(int64(sleepTime)*2, int64(maxSleepTime)))
			fmt.Printf("Sleep time doubled to %v\n", sleepTime)
		}

		time.Sleep(sleepTime)
	}
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

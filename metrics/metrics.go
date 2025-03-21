package metrics

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// runCommand executes a shell command and returns the output as a string.
func runCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return strings.TrimSpace(string(output))
}

// GetMetrics retrieves CPU and Memory usage and includes device ID.
func GetMetrics(deviceID string) map[string]string {
	commands := map[string][]string{
		"cpu":    {"sh", "-c", "grep 'cpu ' /proc/stat | awk '{usage=100-($5*100/($2+$3+$4+$5+$6+$7+$8))} END {print usage}'"},
		"memory": {"sh", "-c", "free | awk '/Mem:/ {print ($3/$2)*100}'"},
	}

	metrics := make(map[string]string)
	metrics["deviceId"] = deviceID                                  // Add device ID to metrics
	metrics["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10) // Unix timestamp

	for name, args := range commands {
		metrics[name] = runCommand(args[0], args[1:]...)
	}

	return metrics
}

// HasSignificantChange checks if there is a significant change in metrics.
func HasSignificantChange(old, new map[string]string, threshold float64) bool {
	for key, newValue := range new {
		// Ignore device ID and Timestamp in comparison
		if key == "deviceId" || key == "timestamp" {
			continue
		}

		oldValue, exists := old[key]
		if !exists {
			return true
		}

		oldFloat, err1 := strconv.ParseFloat(strings.TrimSuffix(oldValue, "%"), 64)
		newFloat, err2 := strconv.ParseFloat(strings.TrimSuffix(newValue, "%"), 64)

		if err1 != nil || err2 != nil {
			return true // Assume change is significant if conversion fails
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

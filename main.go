package main

import (
	"fmt"
	"os/exec"
	// "github.com/gorilla/websocket"
)

func runCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return string(output)
}

func main() {
	commands := map[string][]string{
		"CPU Usage":    {"sh", "-c", "grep 'cpu ' /proc/stat | awk '{usage=100-($5*100/($2+$3+$4+$5+$6+$7+$8))} END {print usage \"%\"}'"},
		"Memory Usage": {"sh", "-c", "free | awk '/Mem:/ {print ($3/$2)*100 \"%\"}'"},
		"Disk Usage":   {"sh", "-c", "df --output=pcent / | tail -n1 | tr -d ' \"%\"'"},
	}

	// Connect to WebSocket server
	// wsURL := "ws://backend-server.com/metrics"
	// conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	// if err != nil {
	// 	log.Fatal("WebSocket connection failed:", err)
	// }
	// defer conn.Clos

	// Collect and send metrics
	for name, args := range commands {
		result := runCommand(args[0], args[1:]...)
		message := fmt.Sprintf("%s: %s", name, result)
		// if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		// 	log.Println("Error sending message:", err)
		// }
		fmt.Printf("%V", message)
	}
}

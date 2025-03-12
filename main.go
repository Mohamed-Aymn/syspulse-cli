package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"syspulse-cli/config"
	"syspulse-cli/identifiers"
	"syspulse-cli/metrics"
	"syspulse-cli/websocket"
)

const pidFile = "/tmp/syspulse-cli.pid"

func main() {
	// Define flags
	daemon := flag.Bool("d", false, "Run in background mode")
	key := flag.String("key", "", "Provide the key generated from the website")
	stop := flag.Bool("stop", false, "Stop the running daemon")
	flag.Parse()

	// Stop the process if -stop is provided
	if *stop {
		stopProcess()
		return
	}

	// Ensure key is provided for starting the process
	if *key == "" {
		fmt.Println("Error: -key flag is required")
		os.Exit(1)
	}

	if !*daemon {
		// Relaunch in background
		cmd := exec.Command(os.Args[0], "-d", "-key", *key) // Re-run with -d flag
		cmd.Stdout = nil                                    // Hide output
		cmd.Stderr = nil                                    // Hide errors
		cmd.Stdin = nil                                     // Detach from terminal
		cmd.Start()                                         // Start in the background

		// Save the process ID
		savePID(cmd.Process.Pid)
		fmt.Println("Process started in the background. PID:", cmd.Process.Pid)
		os.Exit(0) // Exit parent process
	}

	// Background process starts here
	savePID(os.Getpid()) // Save PID for stopping later
	runApplication(*key)

	// Cleanup PID file on exit
	os.Remove(pidFile)
}

// Saves the process ID to a file
func savePID(pid int) {
	f, err := os.Create(pidFile)
	if err != nil {
		fmt.Println("Failed to create PID file:", err)
		os.Exit(1)
	}
	defer f.Close()
	fmt.Fprintln(f, pid)
}

// Stops the running daemon process
func stopProcess() {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Println("No running process found.")
		return
	}

	// Trim any whitespace (including newline)
	pidStr := string(data)
	pidStr = strings.TrimSpace(pidStr)

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		fmt.Println("Invalid PID file:", err)
		return
	}

	// Send SIGTERM to the process
	err = syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		fmt.Println("Failed to stop process:", err)
	} else {
		fmt.Println("Process stopped successfully.")
		os.Remove(pidFile)
	}
}

func runApplication(key string) {
	isNew := !identifiers.IsIDStored()

	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	resp, err := identifiers.ValidateKey(key, isNew)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if isNew {
		id, err := identifiers.FetchAndStoreID(resp)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("ID saved successfully:", id)
	}

	deviceID, err := identifiers.ReadIDFromFile()
	if err != nil {
		fmt.Println("Error retrieving device ID:", err)
		os.Exit(1)
	}

	u := url.URL{Scheme: "ws", Host: cfg.Server.URL, Path: "/", RawQuery: "deviceId=" + deviceID}
	fmt.Printf("Connecting to %s\n", u.String())

	client, err := websocket.NewWebSocketClient(u.String())
	if err != nil {
		fmt.Println("WebSocket connection failed:", err)
		return
	}
	defer client.Close()

	lastMetrics := make(map[string]string)
	sleepTime := 5 * time.Second
	skipCounter := 0
	skipLimit := 12

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

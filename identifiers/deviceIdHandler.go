package identifiers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func SaveIDToFile(id string) error {
	const dirPath = "/etc/syspulse-cli"
	const filePath = "/etc/syspulse-cli/id"

	// Ensure the directory exists
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
	}

	// Write the ID to the file
	if err := os.WriteFile(filePath, []byte(id), 0644); err != nil {
		return fmt.Errorf("failed to write ID to %s: %v", filePath, err)
	}

	return nil
}

func ReadIDFromFile() (string, error) {
	const filePath = "/etc/syspulse-cli/id"

	// Read the ID from the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // File does not exist, return empty ID
		}
		return "", fmt.Errorf("failed to read ID from %s: %v", filePath, err)
	}

	return string(data), nil
}

func IsIDStored() bool {
	const filePath = "/etc/syspulse-cli/id"

	// Check if the file exists
	_, err := os.Stat(filePath)
	return err == nil || !os.IsNotExist(err) // Returns true if file exists, false otherwise
}

func GetIDFromResponse(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", fmt.Errorf("error decoding JSON response: %v", err)
	}

	id, ok := responseData["deviceId"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format: 'id' not found")
	}

	return id, nil
}

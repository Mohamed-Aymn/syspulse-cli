package identifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"syspulse-cli/config"
)

// ValidateKey sends a key validation request and returns the response for further processing.
func ValidateKey(key string, isNew bool) (*http.Response, error) {
	// Load config (singleton, already cached after first call)
	envCfg := config.ReadEnvConfig() // Load env config (singleton)
	cfg, err := config.ReadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	// Determine protocol based on environment
	protocol := "https"
	url := cfg.Server.URL
	if envCfg.ENV == "development" {
		protocol = "http"
		url = "localhost:3000"
	}

	// Construct URL from config
	apiURL := fmt.Sprintf("%s://%s/api/keys/validate", protocol, url)

	// Create request body
	data := map[string]interface{}{
		"key":   key,
		"isNew": isNew,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON: %w", err)
	}

	// Send the POST request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close() // Close response body on error
		return nil, fmt.Errorf("server responded with: %s", resp.Status)
	}

	return resp, nil
}

// FetchAndStoreID extracts and stores the ID if the key is new
func FetchAndStoreID(resp *http.Response) (string, error) {
	id, err := GetIDFromResponse(resp)
	if err != nil {
		return "", err
	}

	if err := SaveIDToFile(id); err != nil {
		return "", fmt.Errorf("error saving ID: %w", err)
	}

	return id, nil
}

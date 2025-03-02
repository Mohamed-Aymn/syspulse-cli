package identifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ValidateKey sends a key validation request and returns the response for further processing.
func ValidateKey(key string, isNew bool) (*http.Response, error) {
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
	resp, err := http.Post("http://localhost:3000/api/keys/validate", "application/json", bytes.NewBuffer(jsonData))
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

package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type authorisationResponse struct {
	URL string `json:"url"`
}

func GenerateConnectionURL(handle, token string) (string, error) {
	encodedHandle := url.QueryEscape("@" + handle)
	endpoint := fmt.Sprintf("https://discord.com/api/v9/connections/mastodon/authorize?handle=%s", encodedHandle)

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("the HTTP request could not be constructed: %w", err)
	}
	req.Header.Set("authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("the request to the Discord API endpoint was unsuccessful: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("the response body could not be read: %w", err)
	}

	var result authorisationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("the Discord API response could not be parsed: %w", err)
	}

	if result.URL == "" {
		return "", errors.New("the Discord API did not return a valid authorisation URL; verify that the supplied token is correct and that network connectivity is available")
	}

	return result.URL, nil
}

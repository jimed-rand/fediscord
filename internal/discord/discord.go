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

type authResponse struct {
	URL string `json:"url"`
}

func GenerateConnectionURL(handle, token string) (string, error) {
	encodedHandle := url.QueryEscape("@" + handle)
	endpoint := fmt.Sprintf("https://discord.com/api/v9/connections/mastodon/authorize?handle=%s", encodedHandle)

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request to Discord API failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var result authResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse Discord response: %w", err)
	}

	if result.URL == "" {
		return "", errors.New("Discord did not return an authorization URL; check your token or network connection")
	}

	return result.URL, nil
}

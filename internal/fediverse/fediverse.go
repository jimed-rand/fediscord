package fediverse

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var handleRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

var incompatiblePlatforms = []string{"misskey", "firefish", "calckey", "foundkey"}

type InstanceInfo struct {
	Version string `json:"version"`
}

func ValidateHandle(handle string) (string, error) {
	handle = strings.TrimPrefix(handle, "@")
	if !handleRegex.MatchString(handle) {
		return "", errors.New("invalid Fediverse handle format; expected: username@instance.domain")
	}
	return handle, nil
}

func ExtractInstance(handle string) string {
	parts := strings.SplitN(handle, "@", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func CheckMastodonAPISupport(instance string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("https://%s/api/v1/instance", instance)

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("could not reach instance: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var info InstanceInfo
	if err := json.Unmarshal(body, &info); err != nil || info.Version == "" {
		return "", errors.New("instance did not return a valid Mastodon API response")
	}

	versionLower := strings.ToLower(info.Version)
	for _, platform := range incompatiblePlatforms {
		if strings.Contains(versionLower, platform) {
			return info.Version, fmt.Errorf("instance is running %s which does not support Mastodon API connections", info.Version)
		}
	}

	return info.Version, nil
}

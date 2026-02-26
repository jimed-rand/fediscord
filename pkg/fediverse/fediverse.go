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
		return "", errors.New("the supplied Fediverse handle does not conform to the expected format; the required format is: username@instance.domain")
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
	apiURL := fmt.Sprintf("https://%s/api/v1/instance", instance)

	resp, err := client.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("the specified instance could not be reached: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("the instance response could not be read: %w", err)
	}

	var info InstanceInfo
	if err := json.Unmarshal(body, &info); err != nil || info.Version == "" {
		return "", errors.New("the instance did not return a valid Mastodon API v1 response; it may not be compatible with this tool")
	}

	versionLower := strings.ToLower(info.Version)
	for _, platform := range incompatiblePlatforms {
		if strings.Contains(versionLower, platform) {
			return info.Version, fmt.Errorf("the instance is operating on %s, which does not implement the Mastodon API and is therefore incompatible with this tool", info.Version)
		}
	}

	return info.Version, nil
}

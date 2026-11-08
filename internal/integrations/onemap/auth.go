package onemap

import (
	"bytes"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

const endpoint = "https://www.onemap.gov.sg/api/auth/post/getToken"

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	AccessToken string `json:"access_token"`
	Expiry      int64  `json:"expiry_timestamp,string"`
}

type cachedToken struct {
	sync.Mutex
	accessToken string
	expiresAt   int64
}

var tokenCache cachedToken

// AccessToken returns a cached token, authenticating when it is absent or expired.
func AccessToken() (string, error) {
	tokenCache.Lock()
	defer tokenCache.Unlock()

	if tokenCache.accessToken != "" && time.Now().Unix() < tokenCache.expiresAt {
		return tokenCache.accessToken, nil
	}

	payload, err := json.Marshal(authRequest{
		Email:    os.Getenv("ONE_MAP_EMAIL"),
		Password: os.Getenv("ONE_MAP_PASSWORD"),
	})
	if err != nil {
		return "", fmt.Errorf("encode OneMap credentials: %w", err)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("authenticate with OneMap: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authenticate with OneMap: HTTP %s", resp.Status)
	}

	var auth authResponse
	if err := json.UnmarshalRead(resp.Body, &auth); err != nil {
		return "", fmt.Errorf("decode OneMap authentication response: %w", err)
	}
	expiresAt := auth.Expiry
	if auth.AccessToken == "" || expiresAt <= time.Now().Unix() {
		return "", fmt.Errorf("invalid OneMap authentication response")
	}
	tokenCache.accessToken, tokenCache.expiresAt = auth.AccessToken, expiresAt
	return tokenCache.accessToken, nil
}

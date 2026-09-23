package core_igdb_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type TokenManager struct {
	client       *http.Client
	clientID     string
	clientSecret string

	mu          sync.Mutex
	accessToken string
	expiresAt   time.Time
}

func NewTokenManager(client *http.Client, clientID string, clientSecret string) *TokenManager {
	return &TokenManager{
		client: client,
		clientID: clientID,
		clientSecret: clientSecret,
	}
}

func (m *TokenManager) GetAccessToken(ctx context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.accessToken != "" && time.Now().Before(m.expiresAt.Add(-5*time.Minute)) {
		return m.accessToken, nil
	}

	token, expiresIn, err := m.fetchToken(ctx)
	if err != nil {
		return "", err
	}

	m.accessToken = token
	m.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	return token, nil
}

func (m *TokenManager) fetchToken(ctx context.Context) (string, int, error) {
	url := fmt.Sprintf(
		"https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials",
		m.clientID, m.clientSecret,
	)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return "", 0, fmt.Errorf("build token request: %w", err)
	}

	response, err := m.client.Do(request)
	if err != nil {
		return "", 0, fmt.Errorf("request token: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return "", 0, fmt.Errorf("token status %d: %s", response.StatusCode, string(body))
	}

	var tokenRespone struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(response.Body).Decode(&tokenRespone); err != nil {
		return "", 0, fmt.Errorf("decode token: %w", err)
	}

	return tokenRespone.AccessToken, tokenRespone.ExpiresIn, nil
}

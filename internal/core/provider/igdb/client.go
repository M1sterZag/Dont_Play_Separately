package core_igdb_provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	clientID   string
	tokens     *TokenManager
	limiter    *rate.Limiter
	semaphore  chan struct{}
}

func NewClient(config Config) *Client {
	httpClient := &http.Client{Timeout: config.Timeout}

	return &Client{
		baseURL:    config.BaseURL,
		httpClient: httpClient,
		clientID:   config.ClientID,
		tokens:     NewTokenManager(httpClient, config.ClientID, config.ClientSecret),
		limiter:    rate.NewLimiter(rate.Limit(4), 1),
		semaphore:  make(chan struct{}, 8),
	}
}

func (c *Client) do(ctx context.Context, endpoint, body string) ([]byte, error) {
	select {
	case c.semaphore <- struct{}{}:
		defer func() { <-c.semaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if err := c.limiter.Wait(ctx); err != nil {
		return nil, err
	}

	token, err := c.tokens.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost,
			c.baseURL+"/"+endpoint, strings.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("build igdb request: %w", err)
		}
		req.Header.Set("Client-ID", c.clientID)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("do request: %w", err)
		}

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			respBody, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("igdb status %d: %s", resp.StatusCode, string(respBody))

			backoff := time.Duration(1<<attempt) * time.Second // 1s, 2s, 4s
			select {
			case <-time.After(backoff):
				continue
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
		defer resp.Body.Close()

		return io.ReadAll(resp.Body)
	}
	return nil, lastErr
}

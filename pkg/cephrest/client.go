package cephrest

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client struct {
	apiURL   string
	username string
	password string
	client   *http.Client
}

func NewClient(apiURL, username, password string) *Client {
	transport := &http.Transport{ //nolint:exhaustruct
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec,exhaustruct
	}
	httpClient := &http.Client{ //nolint:exhaustruct
		Transport: transport,
	}

	return &Client{
		apiURL:   apiURL,
		username: username,
		password: password,
		client:   httpClient,
	}
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (c *Client) Auth(ctx context.Context) (*AuthResponse, error) {
	url := c.apiURL + "/api/auth"

	body, err := json.Marshal(map[string]string{
		"username": c.username,
		"password": c.password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.ceph.api.v1.0+json")

	code, content, err := doRequest(c.client, req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if code != http.StatusCreated {
		return nil, fmt.Errorf("%w: %d: %s", ErrUnexpectedStatusCode, code, string(content))
	}

	var authResponse AuthResponse

	err = json.Unmarshal(content, &authResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &authResponse, nil
}

func (c *Client) Logout(ctx context.Context, token string) error {
	url := c.apiURL + "/api/auth/logout"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.ceph.api.v1.0+json")

	code, content, err := doRequest(c.client, req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	if code != http.StatusOK {
		return fmt.Errorf("%w: %d: %s", ErrUnexpectedStatusCode, code, string(content))
	}

	return nil
}

func (c *Client) GetHealthFull(ctx context.Context, token string) error {
	url := c.apiURL + "/api/health/full"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.ceph.api.v1.0+json")

	code, content, err := doRequest(c.client, req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	if code != http.StatusOK {
		return fmt.Errorf("%w: %d: %s", ErrUnexpectedStatusCode, code, string(content))
	}

	log.Println("cluster", string(content))

	return nil
}

func doRequest(client *http.Client, req *http.Request) (int, []byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to execute request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp.StatusCode, content, nil
}

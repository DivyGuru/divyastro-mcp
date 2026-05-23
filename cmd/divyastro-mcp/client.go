package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// apiClient is a thin wrapper around http.Client that knows how to
// talk to api.divyastroapi.com — base URL, bearer auth, JSON decode,
// useful error messages.
//
// Connection-pooled by design: multiple goroutines from the MCP server
// can call concurrent tools without thrashing TCP. The default
// http.Client transport already pools, but we tune the timeout to a
// reasonable per-request limit (15s — most endpoints respond in
// under 100ms; PDF endpoints take longer but they're not exposed in
// MCP v1).
type apiClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	userAgent  string
}

func newAPIClient(baseURL, apiKey, version string) *apiClient {
	return &apiClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		userAgent: "divyastro-mcp/" + version,
	}
}

// get performs an authenticated GET against the upstream API and
// decodes the JSON response into out. Query parameters are passed as
// a url.Values map so callers don't have to construct query strings
// manually.
//
// On 4xx/5xx, the upstream JSON error body (if any) is surfaced in
// the returned error so the AI client sees a useful message instead
// of "request failed". This is critical for debugging — when an AI
// assistant says "the tool failed", users want to know why.
func (c *apiClient) get(ctx context.Context, path string, params url.Values, out any) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1 MiB cap
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		// Try to surface the structured error message — the upstream
		// API returns {"error": "...", "code": "...", "details": "..."}.
		// Falling back to raw body keeps unknown formats debuggable.
		var apiErr struct {
			Error   string `json:"error"`
			Code    string `json:"code"`
			Details string `json:"details"`
		}
		if jerr := json.Unmarshal(body, &apiErr); jerr == nil && apiErr.Error != "" {
			msg := apiErr.Error
			if apiErr.Details != "" {
				msg += " — " + apiErr.Details
			}
			return fmt.Errorf("API %d %s: %s", resp.StatusCode, apiErr.Code, msg)
		}
		return fmt.Errorf("API %d: %s", resp.StatusCode, string(body))
	}

	if out == nil {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

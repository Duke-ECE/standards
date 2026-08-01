// Package postgrest implements the slice ports against Supabase PostgREST.
// Tables it touches must have RLS enabled with no policies: only the
// service key can read or write them, and only this service holds it.
package postgrest

import (
	"context"
	"fmt"
	"net/http"
)

// Client is the shared PostgREST client: base URL, service key, and the
// auth headers every request needs.
type Client struct {
	baseURL    string // e.g. https://<ref>.supabase.co
	serviceKey string
	http       *http.Client
}

// NewClient returns a Client. A nil httpClient uses http.DefaultClient.
func NewClient(baseURL, serviceKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{baseURL: baseURL, serviceKey: serviceKey, http: httpClient}
}

// newRequest stamps the service-key auth headers shared by every call.
func (c *Client) newRequest(ctx context.Context, method, path string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+"/rest/v1/"+path, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("apikey", c.serviceKey)
	req.Header.Set("Authorization", "Bearer "+c.serviceKey)
	return req, nil
}

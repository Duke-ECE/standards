package store

import (
	"context"
	"fmt"
	"net/http"
)

// postgrest implements Store against Supabase PostgREST. Tables it touches
// must have RLS enabled with no policies: only the service key can read or
// write them, and only this service holds that key.
type postgrest struct {
	baseURL    string // e.g. https://<ref>.supabase.co
	serviceKey string
	client     *http.Client
}

// NewPostgREST returns a Store backed by Supabase PostgREST.
// A nil httpClient uses http.DefaultClient.
func NewPostgREST(baseURL, serviceKey string, httpClient *http.Client) Store {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &postgrest{baseURL: baseURL, serviceKey: serviceKey, client: httpClient}
}

// newRequest stamps the service-key auth headers shared by every call.
func (p *postgrest) newRequest(ctx context.Context, method, path string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, p.baseURL+"/rest/v1/"+path, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("apikey", p.serviceKey)
	req.Header.Set("Authorization", "Bearer "+p.serviceKey)
	return req, nil
}

func (p *postgrest) Ping(ctx context.Context) error {
	req, err := p.newRequest(ctx, http.MethodGet, "")
	if err != nil {
		return err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("ping postgrest: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("ping postgrest: unexpected status %d", resp.StatusCode)
	}
	return nil
}

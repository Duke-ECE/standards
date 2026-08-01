package postgrest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Duke-ECE/go-service-template/internal/example"
)

// ExampleStore implements example.Store over the "examples" table.
// Behavioral parity with the memory implementation is required: same
// errors, same semantics.
type ExampleStore struct {
	c *Client
}

func NewExampleStore(c *Client) *ExampleStore {
	return &ExampleStore{c: c}
}

// do sends a request and decodes the JSON array response into rows.
func (c *Client) do(ctx context.Context, method, path string, body any, rows any) error {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode body: %w", err)
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+"/rest/v1/"+path, rdr)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("apikey", c.serviceKey)
	req.Header.Set("Authorization", "Bearer "+c.serviceKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("postgrest %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("postgrest %s %s: unexpected status %d", method, path, resp.StatusCode)
	}
	if rows != nil {
		if err := json.NewDecoder(resp.Body).Decode(rows); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

func (s *ExampleStore) Create(ctx context.Context, e *example.Example) error {
	return s.c.do(ctx, "POST", "examples", e, nil)
}

func (s *ExampleStore) Get(ctx context.Context, id string) (*example.Example, error) {
	var rows []example.Example
	if err := s.c.do(ctx, "GET", "examples?id=eq."+url.QueryEscape(id), nil, &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, example.ErrNotFound // domain error, not a raw PostgREST one
	}
	return &rows[0], nil
}

var _ example.Store = (*ExampleStore)(nil)

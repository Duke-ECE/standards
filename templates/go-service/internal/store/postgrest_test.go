package store

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Fake external APIs with httptest: assert the request shape, script the
// response. No docker, no testcontainers.
func TestPostgRESTPingSendsServiceKey(t *testing.T) {
	var gotAPIKey, gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAPIKey = r.Header.Get("apikey")
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(ts.Close)

	s := NewPostgREST(ts.URL, "test-key", nil)
	if err := s.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
	if gotAPIKey != "test-key" || gotAuth != "Bearer test-key" {
		t.Errorf("auth headers = %q / %q", gotAPIKey, gotAuth)
	}
}

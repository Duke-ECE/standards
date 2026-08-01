package postgrest_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Duke-ECE/go-service-template/internal/example"
	"github.com/Duke-ECE/go-service-template/internal/infrastructure/postgrest"
)

// Fake external APIs with httptest: assert the request shape, script the
// response. No docker, no testcontainers.
func TestGetSendsServiceKey(t *testing.T) {
	var gotAPIKey, gotAuth, gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAPIKey = r.Header.Get("apikey")
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.RequestURI()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"ex-1","name":"ada"}]`))
	}))
	t.Cleanup(ts.Close)

	s := postgrest.NewExampleStore(postgrest.NewClient(ts.URL, "test-key", nil))
	e, err := s.Get(context.Background(), "ex-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if e.ID != "ex-1" {
		t.Fatalf("unexpected example: %+v", e)
	}
	if gotAPIKey != "test-key" || gotAuth != "Bearer test-key" {
		t.Errorf("auth headers = %q / %q", gotAPIKey, gotAuth)
	}
	if gotPath != "/rest/v1/examples?id=eq.ex-1" {
		t.Errorf("path = %q", gotPath)
	}
}

// Empty results are domain errors, not driver errors.
func TestGetNotFoundMapsToDomainError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	t.Cleanup(ts.Close)

	s := postgrest.NewExampleStore(postgrest.NewClient(ts.URL, "test-key", nil))
	if _, err := s.Get(context.Background(), "ex-nope"); err != example.ErrNotFound {
		t.Fatalf("err = %v, want example.ErrNotFound", err)
	}
}

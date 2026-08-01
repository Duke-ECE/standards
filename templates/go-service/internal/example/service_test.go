package example_test

import (
	"context"
	"testing"

	"github.com/Duke-ECE/go-service-template/internal/example"
	"github.com/Duke-ECE/go-service-template/internal/infrastructure/memory"
)

// Slice tests run the real service against the memory store — no servers,
// no fakes to maintain. The memory backend is the reference implementation.
func TestCreateAndGet(t *testing.T) {
	svc := example.NewService(memory.NewExampleStore())

	e, err := svc.Create(context.Background(), "ada")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	got, err := svc.Get(context.Background(), e.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != "ada" {
		t.Fatalf("Name = %q", got.Name)
	}
}

func TestCreateEmptyNameIsDomainError(t *testing.T) {
	svc := example.NewService(memory.NewExampleStore())
	if _, err := svc.Create(context.Background(), ""); err != example.ErrEmptyName {
		t.Fatalf("err = %v, want example.ErrEmptyName", err)
	}
}

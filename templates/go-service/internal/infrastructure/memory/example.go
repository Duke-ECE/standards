// Package memory is the in-memory implementation of the slice ports.
// It is REQUIRED for every port: the binary must run locally with zero
// external services, and slice tests run against this implementation.
package memory

import (
	"context"
	"sync"

	"github.com/Duke-ECE/go-service-template/internal/example"
)

// ExampleStore implements example.Store in memory.
type ExampleStore struct {
	mu sync.RWMutex
	by map[string]*example.Example
}

func NewExampleStore() *ExampleStore {
	return &ExampleStore{by: make(map[string]*example.Example)}
}

func (s *ExampleStore) Create(_ context.Context, e *example.Example) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := *e // copy in, copy out: callers can't mutate stored state
	s.by[e.ID] = &cp
	return nil
}

func (s *ExampleStore) Get(_ context.Context, id string) (*example.Example, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.by[id]
	if !ok {
		return nil, example.ErrNotFound
	}
	cp := *e
	return &cp, nil
}

var _ example.Store = (*ExampleStore)(nil)

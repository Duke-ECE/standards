package example

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"
)

// Service holds the slice's business rules. Transports call it; it never
// calls them. It depends on the Store port, not on any implementation.
type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

// Create demonstrates a rule-bearing method: ids are generated here (not
// in the store), and duplicates are domain errors, not driver errors.
func (s *Service) Create(ctx context.Context, name string) (*Example, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return nil, err
	}
	e := &Example{ID: "ex-" + hex.EncodeToString(b[:]), Name: name, CreatedAt: time.Now().UTC()}
	if err := s.store.Create(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Example, error) {
	return s.store.Get(ctx, id)
}

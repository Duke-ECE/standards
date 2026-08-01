package example

import "context"

// Store is the persistence port, OWNED BY THE SLICE: it is declared here
// (where it is consumed), never next to its implementations. Methods are
// domain-shaped, not SQL-shaped. infrastructure/* implements it.
type Store interface {
	Create(ctx context.Context, e *Example) error
	Get(ctx context.Context, id string) (*Example, error)
}

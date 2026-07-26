// Package store is the persistence seam. Business logic depends on the
// Store interface; tests substitute an in-memory fake. The PostgREST
// implementation is the standard way services reach Supabase — delete this
// package entirely if your service is stateless.
package store

import (
	"context"
)

// Store is the persistence contract for the service's domain. Keep methods
// domain-shaped (what the logic needs), not SQL-shaped.
type Store interface {
	// Example; replace with real domain methods.
	Ping(ctx context.Context) error
}

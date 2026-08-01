// Package example is a domain slice — the unit of vertical organization.
// A slice owns everything about its aggregate in one package: types,
// rules (service.go), the storage port (store.go), and domain errors
// (errors.go). It imports nothing from transport or infrastructure.
// Rename it; every service needs at least one real slice.
package example

import "time"

// Example is the slice's aggregate type.
type Example struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

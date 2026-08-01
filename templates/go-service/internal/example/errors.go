package example

import "errors"

// Domain errors. Transports map these to gRPC statuses / HTTP codes in
// exactly one place; infrastructure must translate driver errors into
// these rather than leaking them upward.
var (
	ErrNotFound  = errors.New("example: not found")
	ErrEmptyName = errors.New("example: name is required")
)

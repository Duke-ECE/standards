// Package server constructs the gRPC server and registers services.
// Construction lives here; business logic lives in internal/<domain>;
// persistence lives behind the internal/store interface.
package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// New builds the gRPC server with all services registered.
func New() *grpc.Server {
	s := grpc.NewServer()

	// Register generated services here, e.g.:
	//   v1.RegisterEchoServiceServer(s, echo.NewServer(store))
	// For HTTP (Gin) services instead, see the standards AGENTS.md §2:
	// construct an explicit http.Server with ReadHeaderTimeout and register
	// routes in a separate routes.go (RegisterRoutes).

	// Reflection enables grpcurl debugging; safe to keep in production.
	reflection.Register(s)
	return s
}

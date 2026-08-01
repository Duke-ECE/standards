// Package grpc constructs the gRPC server and registers services.
// Handlers live here (thin: decode → call slice → map error); business
// rules live in the slices; construction is wired from cmd/server.
package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// New builds the gRPC server with all services registered.
func New() *grpc.Server {
	s := grpc.NewServer()

	// Register generated services here, e.g.:
	//   v1.RegisterExampleServiceServer(s, NewExampleHandler(exampleSvc))
	// Handlers are thin: they decode the request, call ONE slice method,
	// and map the domain error in errors.go — the only mapping place.

	// Reflection enables grpcurl debugging; safe to keep in production.
	reflection.Register(s)
	return s
}

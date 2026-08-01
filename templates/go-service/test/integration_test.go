// Package test holds whole-service integration tests: the service is
// assembled exactly as cmd/server/main.go assembles it and driven through
// its public API over a real connection. Fakes belong only at true
// external boundaries — the memory backends are real implementations, so
// no fakes are needed here at all.
package test

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"

	"github.com/Duke-ECE/go-service-template/internal/example"
	"github.com/Duke-ECE/go-service-template/internal/infrastructure/memory"
	transportgrpc "github.com/Duke-ECE/go-service-template/internal/transport/grpc"
)

// Boots the full stack on a real port, like main does.
func TestServerBootsAndServesReflection(t *testing.T) {
	// Same wiring as main (memory backends → slices → transports).
	_ = example.NewService(memory.NewExampleStore())
	srv := transportgrpc.New()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	stream, err := grpc_reflection_v1.NewServerReflectionClient(conn).ServerReflectionInfo(context.Background())
	if err != nil {
		t.Fatalf("reflection: %v", err)
	}
	if err := stream.Send(&grpc_reflection_v1.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1.ServerReflectionRequest_ListServices{},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := stream.Recv(); err != nil {
		t.Fatalf("reflection recv: %v", err)
	}
}

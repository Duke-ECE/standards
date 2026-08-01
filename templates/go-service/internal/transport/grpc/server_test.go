package grpc

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/grpc/test/bufconn"
)

// The bufconn pattern: serve in-process and dial through a custom dialer —
// no ports, no network, fast and hermetic. All gRPC tests use this.
func newTestConn(t *testing.T) *grpc.ClientConn {
	t.Helper()
	lis := bufconn.Listen(1024 * 1024)
	s := New()
	go func() { _ = s.Serve(lis) }()
	t.Cleanup(s.Stop)

	conn, err := grpc.NewClient("passthrough:///bufconn",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial bufconn: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

// Reflection being served proves the server booted and registered cleanly.
func TestServerServesReflection(t *testing.T) {
	conn := newTestConn(t)
	stream, err := grpc_reflection_v1.NewServerReflectionClient(conn).
		ServerReflectionInfo(context.Background())
	if err != nil {
		t.Fatalf("reflection info: %v", err)
	}
	if err := stream.Send(&grpc_reflection_v1.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1.ServerReflectionRequest_ListServices{},
	}); err != nil {
		t.Fatalf("reflection send: %v", err)
	}
	resp, err := stream.Recv()
	if err != nil {
		t.Fatalf("reflection recv: %v", err)
	}
	if resp.GetListServicesResponse() == nil {
		t.Fatal("expected a service list from reflection")
	}
}

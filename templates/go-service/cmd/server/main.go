// Command server is the ONLY assembly point of the app:
// config → infrastructure → domain services → transports → serve.
// Everything below main is injectable and testable; nothing below main
// knows how the pieces are wired.
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	transportgrpc "github.com/Duke-ECE/go-service-template/internal/transport/grpc"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	// --- infrastructure → slices → transports ---
	// Memory backends are the default so the binary runs with zero
	// external services; swap in postgrest.NewClient(...) implementations
	// when SUPABASE_URL + SUPABASE_SERVICE_ROLE_KEY are set.
	//   store := memory.NewExampleStore()
	//   exampleSvc := example.NewService(store)
	grpcServer := transportgrpc.New()

	// GracefulStop on SIGINT/SIGTERM: in-flight RPCs finish, then we exit.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		log.Println("shutting down gracefully")
		grpcServer.GracefulStop()
	}()

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listen :%s: %v", port, err)
	}
	log.Printf("listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

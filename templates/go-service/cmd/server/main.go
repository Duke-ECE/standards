// Command server boots the service: env parsing, wiring, and graceful
// shutdown only — no business logic belongs here.
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Duke-ECE/go-service-template/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listen :%s: %v", port, err)
	}

	grpcServer := server.New()

	// GracefulStop on SIGINT/SIGTERM: in-flight RPCs finish, then we exit.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		log.Println("shutting down gracefully")
		grpcServer.GracefulStop()
	}()

	log.Printf("listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

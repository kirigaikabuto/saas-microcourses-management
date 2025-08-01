package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kirigaikabuto/saas-microcourses-management/internal/service"
	companyv1 "github.com/kirigaikabuto/saas-microcourses-management/proto/gen/proto/company/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	// Database connection
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://admin:password@postgres:5432/saas_microcourses?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}
	defer pool.Close()

	// Test database connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create gRPC server
	server := grpc.NewServer()

	// Register services
	companyService := service.NewCompanyService(pool)
	companyv1.RegisterCompanyServiceServer(server, companyService)

	// Enable reflection for grpcurl/grpc clients
	reflection.Register(server)

	// Listen on port
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on port %s", port)

	// Handle graceful shutdown
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gRPC server...")
	server.GracefulStop()
}

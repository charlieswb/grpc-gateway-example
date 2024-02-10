package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	pb "restgrpc/pkg/protobuf"
	"restgrpc/pkg/service"
)

func main() {
	gs := startGRPCServer()
	hs := startHTTPServer()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)

	sig := <-exit
	slog.Warn("Caught interrupt signal", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if hs != nil {
		slog.Info("Shutting down http server")
		hs.Shutdown(ctx)
	}
	if gs != nil {
		slog.Info("Shutting down grpc server")
		gs.GracefulStop()
	}

	slog.Info("Running cleanup tasks...")
}

func startHTTPServer() *http.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 8000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterQuoteServiceHandlerFromEndpoint(context.Background(), mux, "localhost:9000", opts); err != nil {
		log.Fatalf("failed to register HTTP quote service handler: %v", err)
	}

	httpServer := &http.Server{
		Handler: mux,
	}

	go func() {
		slog.Info("Listening to HTTP request at port 8000")
		if err := httpServer.Serve(lis); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to create HTTP server: %v", err)
		}
	}()

	return httpServer
}

func startGRPCServer() *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// register gRPC service
	pb.RegisterQuoteServiceServer(grpcServer, service.NewQuoteService())
	reflection.Register(grpcServer)

	go func() {
		slog.Info("Listening to gRPC request at port 9000")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to create gRPC server: %v", err)
		}
	}()

	return grpcServer
}

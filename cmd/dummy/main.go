package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	grpcproxy "github.com/rollkit/go-execution/proxy/grpc"
	"github.com/rollkit/go-execution/test"
	pb "github.com/rollkit/go-execution/types/pb/execution"
)

func main() {
	listenAddress, err := parseListenAddress()
	if err != nil {
		log.Fatalf("Failed to parse listen address: %v\n", err)
	}
	listener, err := net.Listen("tcp4", listenAddress)
	if err != nil {
		log.Fatalf("Failed to listen on %q: %v\n", listenAddress, err)
	}
	defer func() {
		_ = listener.Close()
	}()

	log.Println("Creating Dummy Executor and gRPC server")
	dummy := test.NewDummyExecutor()
	server := grpcproxy.NewServer(dummy, grpcproxy.DefaultConfig())
	s := grpc.NewServer()
	pb.RegisterExecutionServiceServer(s, server)

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	doneChan := make(chan interface{}, 1)
	go func() {
		log.Printf("Serving (%s)...\n", listenAddress)
		log.Println("Type Ctrl+C to shutdown")
		if err := s.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf("Server exited with error: %v\n", err)
		}
		doneChan <- nil
	}()

	// Handle shutdown signal
	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		s.GracefulStop()
	}()

	<-doneChan
	log.Println("Server stopped")
}

func parseListenAddress() (string, error) {
	var listenAddress string
	flag.StringVar(&listenAddress, "address", "127.0.0.1:40041", "gRPC server listen address")
	flag.Parse()

	_, port, err := net.SplitHostPort(listenAddress)
	if err != nil {
		return "", fmt.Errorf("invalid address format %q: %v", listenAddress, err)
	}
	if port == "" {
		return "", errors.New("port cannot be empty")
	}

	return listenAddress, nil
}

package main

import (
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcproxy "github.com/rollkit/go-execution/proxy/grpc"
	"github.com/rollkit/go-execution/test"
	pb "github.com/rollkit/go-execution/types/pb/execution"
)

const bufSize = 1024 * 1024

func main() {
	dummy := test.NewDummyExecutor()

	listenAddress := "127.0.0.1:40041"
	if len(os.Args) == 2 {
		listenAddress = os.Args[1]
	}

	listener, err := net.Listen("tcp4", listenAddress)
	if err != nil {
		log.Fatalf("error while creating listener: %v\n", err)
	}
	defer func() {
		_ = listener.Close()
	}()

	log.Println("Starting server...")
	server := grpcproxy.NewServer(dummy, grpcproxy.DefaultConfig())
	s := grpc.NewServer()
	pb.RegisterExecutionServiceServer(s, server)

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	doneChan := make(chan interface{}, 1)
	go func() {
		log.Println("Serving...")
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

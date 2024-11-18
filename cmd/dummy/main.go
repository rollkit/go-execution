package main

import (
	"errors"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"

	grpcproxy "github.com/rollkit/go-execution/proxy/grpc"
	"github.com/rollkit/go-execution/test"
	pb "github.com/rollkit/go-execution/types/pb/execution"
)

const bufSize = 1024 * 1024

func main() {
	dummy := test.NewDummyExecutor()
	config := &grpcproxy.Config{
		DefaultTimeout: 5 * time.Second,
		MaxRequestSize: bufSize,
	}

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
	server := grpcproxy.NewServer(dummy, config)
	s := grpc.NewServer()
	pb.RegisterExecutionServiceServer(s, server)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		log.Println("Serving...")
		if err := s.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf("Server exited with error: %v\n", err)
		}
		wg.Done()
	}()
	defer s.Stop()

	wg.Wait()
	log.Println("Server stopped")
}

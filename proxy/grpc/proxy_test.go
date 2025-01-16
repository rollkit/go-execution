package grpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	grpcproxy "github.com/rollkit/go-execution/proxy/grpc"
	"github.com/rollkit/go-execution/test"
	pb "github.com/rollkit/go-execution/types/pb/execution"
)

const bufSize = 1024 * 1024

func dialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

type ProxyTestSuite struct {
	test.ExecutorSuite
	server  *grpc.Server
	client  *grpcproxy.Client
	cleanup func()
}

func (s *ProxyTestSuite) SetupTest() {
	exec := test.NewDummyExecutor()
	config := &grpcproxy.Config{
		DefaultTimeout: time.Second,
		MaxRequestSize: bufSize,
	}
	server := grpcproxy.NewServer(exec, config)

	listener := bufconn.Listen(bufSize)
	s.server = grpc.NewServer()
	pb.RegisterExecutionServiceServer(s.server, server)

	go func() {
		if err := s.server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			s.T().Errorf("Server exited with error: %v", err)
		}
	}()

	client := grpcproxy.NewClient()
	client.SetConfig(config)

	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithContextDialer(dialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := client.Start("passthrough://bufnet", opts...)
	require.NoError(s.T(), err)

	for i := 0; i < 10; i++ {
		if _, err := client.GetTxs(context.TODO()); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	s.client = client
	s.Exec = client
	s.TxInjector = exec
	s.cleanup = func() {
		_ = client.Stop()
		s.server.Stop()
	}
}

func (s *ProxyTestSuite) TearDownTest() {
	s.cleanup()
}

func TestProxySuite(t *testing.T) {
	suite.Run(t, new(ProxyTestSuite))
}

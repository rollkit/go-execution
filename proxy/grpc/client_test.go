package grpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/rollkit/rollkit/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/LastL2/go-execution/mocks"
	grpcproxy "github.com/LastL2/go-execution/proxy/grpc"
	pb "github.com/LastL2/go-execution/types"
)

func TestClientServer(t *testing.T) {
	mockExec := mocks.NewMockExecute(t)
	server := grpcproxy.NewServer(mockExec)

	listener := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	pb.RegisterExecutionServiceServer(s, server)

	go func() {
		if err := s.Serve(listener); err != nil {
			t.Errorf("Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := grpcproxy.NewClient(conn)

	t.Run("InitChain", func(t *testing.T) {
		genesisTime := time.Now().UTC()
		initialHeight := uint64(1)
		chainID := "test-chain"
		expectedStateRoot := types.Hash{1, 2, 3}
		expectedMaxBytes := uint64(1000000)

		mockExec.On("InitChain", genesisTime, initialHeight, chainID).
			Return(expectedStateRoot, expectedMaxBytes, nil)

		stateRoot, maxBytes, err := client.InitChain(genesisTime, initialHeight, chainID)

		require.NoError(t, err)
		assert.Equal(t, expectedStateRoot, stateRoot)
		assert.Equal(t, expectedMaxBytes, maxBytes)
		mockExec.AssertExpectations(t)
	})
}

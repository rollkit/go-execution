package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/rollkit/go-execution/mocks"
	grpcproxy "github.com/rollkit/go-execution/proxy/grpc"
	"github.com/rollkit/go-execution/types"
	pb "github.com/rollkit/go-execution/types/pb/execution"
)

func TestClientServer(t *testing.T) {
	mockExec := mocks.NewMockExecutor(t)
	config := &grpcproxy.Config{
		DefaultTimeout: 5 * time.Second,
		MaxRequestSize: bufSize,
	}
	server := grpcproxy.NewServer(mockExec, config)

	listener := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterExecutionServiceServer(s, server)

	go func() {
		if err := s.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("Server exited with error: %v", err)
		}
	}()
	defer s.Stop()

	client := grpcproxy.NewClient()
	client.SetConfig(config)

	err := client.Start("passthrough://bufnet",
		grpc.WithContextDialer(dialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer func() { _ = client.Stop() }()

	mockExec.On("GetTxs", mock.Anything).Return([]types.Tx{}, nil).Maybe()

	t.Run("InitChain", func(t *testing.T) {
		genesisTime := time.Now().UTC().Truncate(time.Second)
		initialHeight := uint64(1)
		chainID := "test-chain"

		// initialize a new Hash with a fixed size
		expectedStateRoot := make([]byte, 32)
		copy(expectedStateRoot, []byte{1, 2, 3})
		var stateRootHash types.Hash
		copy(stateRootHash[:], expectedStateRoot)

		expectedMaxBytes := uint64(1000000)

		// convert time to Unix and back to ensure consistency
		unixTime := genesisTime.Unix()
		expectedTime := time.Unix(unixTime, 0).UTC()

		mockExec.On("InitChain", mock.Anything, expectedTime, initialHeight, chainID).
			Return(stateRootHash, expectedMaxBytes, nil).Once()

		stateRoot, maxBytes, err := client.InitChain(context.TODO(), genesisTime, initialHeight, chainID)

		require.NoError(t, err)
		assert.Equal(t, stateRootHash, stateRoot)
		assert.Equal(t, expectedMaxBytes, maxBytes)
		mockExec.AssertExpectations(t)
	})
}

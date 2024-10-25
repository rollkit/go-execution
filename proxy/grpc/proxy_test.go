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

	"github.com/LastL2/mocks"
	grpcproxy "github.com/LastL2/proxy/grpc"
	pb "github.com/LastL2/types"
)

func TestProxy(t *testing.T) {
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

	t.Run("GetTxs", func(t *testing.T) {
		expectedTxs := []types.Tx{
			[]byte("tx1"),
			[]byte("tx2"),
		}

		mockExec.On("GetTxs").Return(expectedTxs, nil)

		txs, err := client.GetTxs()

		require.NoError(t, err)
		assert.Equal(t, expectedTxs, txs)
		mockExec.AssertExpectations(t)
	})

	t.Run("ExecuteTxs", func(t *testing.T) {
		txs := []types.Tx{[]byte("tx1"), []byte("tx2")}
		blockHeight := uint64(10)
		timestamp := time.Now().UTC()
		prevStateRoot := types.Hash{4, 5, 6}
		expectedStateRoot := types.Hash{7, 8, 9}
		expectedMaxBytes := uint64(2000000)

		mockExec.On("ExecuteTxs", txs, blockHeight, timestamp, prevStateRoot).
			Return(expectedStateRoot, expectedMaxBytes, nil)

		updatedStateRoot, maxBytes, err := client.ExecuteTxs(txs, blockHeight, timestamp, prevStateRoot)

		require.NoError(t, err)
		assert.Equal(t, expectedStateRoot, updatedStateRoot)
		assert.Equal(t, expectedMaxBytes, maxBytes)
		mockExec.AssertExpectations(t)
	})

	t.Run("SetFinal", func(t *testing.T) {
		blockHeight := uint64(10)

		mockExec.On("SetFinal", blockHeight).Return(nil)

		err := client.SetFinal(blockHeight)

		require.NoError(t, err)
		mockExec.AssertExpectations(t)
	})
}

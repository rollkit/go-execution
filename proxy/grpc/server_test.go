package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/rollkit/rollkit/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LastL2/go-execution/mocks"
	grpcproxy "github.com/LastL2/go-execution/proxy/grpc"
	pb "github.com/LastL2/go-execution/types"
)

func TestServer(t *testing.T) {
	mockExec := mocks.NewMockExecute(t)
	server := grpcproxy.NewServer(mockExec)

	t.Run("InitChain", func(t *testing.T) {
		genesisTime := time.Now().UTC()
		initialHeight := uint64(1)
		chainID := "test-chain"
		expectedStateRoot := types.Hash{1, 2, 3}
		expectedMaxBytes := uint64(1000000)

		mockExec.On("InitChain", genesisTime, initialHeight, chainID).
			Return(expectedStateRoot, expectedMaxBytes, nil)

		resp, err := server.InitChain(context.Background(), &pb.InitChainRequest{
			GenesisTime:   genesisTime.Unix(),
			InitialHeight: initialHeight,
			ChainId:       chainID,
		})

		require.NoError(t, err)
		assert.Equal(t, expectedStateRoot[:], resp.StateRoot)
		assert.Equal(t, expectedMaxBytes, resp.MaxBytes)
		mockExec.AssertExpectations(t)
	})

	t.Run("GetTxs", func(t *testing.T) {
		expectedTxs := []types.Tx{
			[]byte("tx1"),
			[]byte("tx2"),
		}

		mockExec.On("GetTxs").Return(expectedTxs, nil)

		resp, err := server.GetTxs(context.Background(), &pb.GetTxsRequest{})

		require.NoError(t, err)
		assert.Equal(t, len(expectedTxs), len(resp.Txs))
		for i, tx := range expectedTxs {
			assert.Equal(t, tx, resp.Txs[i])
		}
		mockExec.AssertExpectations(t)
	})
}

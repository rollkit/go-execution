package jsonrpc_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rollkit/rollkit/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LastL2/go-execution/mocks"
	jsonrpcproxy "github.com/LastL2/go-execution/proxy/jsonrpc"
)

func TestClientServer(t *testing.T) {
	mockExec := mocks.NewMockExecute(t)
	config := &jsonrpcproxy.Config{
		DefaultTimeout: 5 * time.Second,
		MaxRequestSize: 1024 * 1024,
	}
	server := jsonrpcproxy.NewServer(mockExec, config)

	testServer := httptest.NewServer(server)
	defer testServer.Close()

	client := jsonrpcproxy.NewClient()
	client.SetConfig(config)

	err := client.Start(testServer.URL)
	require.NoError(t, err)
	defer client.Stop()

	t.Run("InitChain", func(t *testing.T) {
		genesisTime := time.Now().UTC().Truncate(time.Second)
		initialHeight := uint64(1)
		chainID := "test-chain"

		expectedStateRoot := make([]byte, 32)
		copy(expectedStateRoot, []byte{1, 2, 3})
		var stateRootHash types.Hash
		copy(stateRootHash[:], expectedStateRoot)

		expectedMaxBytes := uint64(1000000)

		// convert time to Unix and back to ensure consistency
		unixTime := genesisTime.Unix()
		expectedTime := time.Unix(unixTime, 0).UTC()

		mockExec.On("InitChain", expectedTime, initialHeight, chainID).
			Return(stateRootHash, expectedMaxBytes, nil).Once()

		stateRoot, maxBytes, err := client.InitChain(genesisTime, initialHeight, chainID)

		require.NoError(t, err)
		assert.Equal(t, stateRootHash, stateRoot)
		assert.Equal(t, expectedMaxBytes, maxBytes)
		mockExec.AssertExpectations(t)
	})

	t.Run("GetTxs", func(t *testing.T) {
		expectedTxs := []types.Tx{[]byte("tx1"), []byte("tx2")}
		mockExec.On("GetTxs").Return(expectedTxs, nil).Once()

		txs, err := client.GetTxs()
		require.NoError(t, err)
		assert.Equal(t, expectedTxs, txs)
		mockExec.AssertExpectations(t)
	})

	t.Run("ExecuteTxs", func(t *testing.T) {
		txs := []types.Tx{[]byte("tx1"), []byte("tx2")}
		blockHeight := uint64(1)
		timestamp := time.Now().UTC().Truncate(time.Second)

		var prevStateRoot types.Hash
		copy(prevStateRoot[:], []byte{1, 2, 3})

		var expectedStateRoot types.Hash
		copy(expectedStateRoot[:], []byte{4, 5, 6})

		expectedMaxBytes := uint64(1000000)

		// convert time to Unix and back to ensure consistency
		unixTime := timestamp.Unix()
		expectedTime := time.Unix(unixTime, 0).UTC()

		mockExec.On("ExecuteTxs", txs, blockHeight, expectedTime, prevStateRoot).
			Return(expectedStateRoot, expectedMaxBytes, nil).Once()

		updatedStateRoot, maxBytes, err := client.ExecuteTxs(txs, blockHeight, timestamp, prevStateRoot)

		require.NoError(t, err)
		assert.Equal(t, expectedStateRoot, updatedStateRoot)
		assert.Equal(t, expectedMaxBytes, maxBytes)
		mockExec.AssertExpectations(t)
	})

	t.Run("SetFinal", func(t *testing.T) {
		blockHeight := uint64(1)
		mockExec.On("SetFinal", blockHeight).Return(nil).Once()

		err := client.SetFinal(blockHeight)
		require.NoError(t, err)
		mockExec.AssertExpectations(t)
	})
}

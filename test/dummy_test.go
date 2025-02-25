package test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/rollkit/go-execution/types"
)

type DummyTestSuite struct {
	ExecutorSuite
}

func (s *DummyTestSuite) SetupTest() {
	dummy := NewDummyExecutor()
	s.Exec = dummy
	s.TxInjector = dummy
}

func TestDummySuite(t *testing.T) {
	suite.Run(t, new(DummyTestSuite))
}

func (s *DummyTestSuite) TestTxRemoval() {
	t := s.T()
	exec := NewDummyExecutor()
	tx1 := types.Tx([]byte{1, 2, 3})
	tx2 := types.Tx([]byte{3, 2, 1})

	exec.InjectTx(tx1)
	exec.InjectTx(tx2)

	// first execution of GetTxs - nothing special
	txs, err := exec.GetTxs(context.Background())
	require.NoError(t, err)
	require.Len(t, txs, 2)
	require.Contains(t, txs, tx1)
	require.Contains(t, txs, tx2)

	// ExecuteTxs was not called, so 2 txs should still be returned
	txs, err = exec.GetTxs(context.Background())
	require.NoError(t, err)
	require.Len(t, txs, 2)
	require.Contains(t, txs, tx1)
	require.Contains(t, txs, tx2)

	dummyStateRoot := []byte("dummy-state-root")
	state, _, err := exec.ExecuteTxs(context.Background(), []types.Tx{tx1}, 1, time.Now(), dummyStateRoot)
	require.NoError(t, err)
	require.NotEmpty(t, state)

	// ExecuteTxs was called, 1 tx remaining in mempool
	txs, err = exec.GetTxs(context.Background())
	require.NoError(t, err)
	require.Len(t, txs, 1)
	require.NotContains(t, txs, tx1)
	require.Contains(t, txs, tx2)
}

func (s *DummyTestSuite) TestExecuteTxsComprehensive() {
	t := s.T()
	tests := []struct {
		name          string
		txs           []types.Tx
		blockHeight   uint64
		timestamp     time.Time
		prevStateRoot types.Hash
		expectedErr   error
	}{
		{
			name:          "valid multiple transactions",
			txs:           []types.Tx{[]byte("tx1"), []byte("tx2"), []byte("tx3")},
			blockHeight:   1,
			timestamp:     time.Now().UTC(),
			prevStateRoot: types.Hash{1, 2, 3},
			expectedErr:   nil,
		},
		{
			name:          "empty state root",
			txs:           []types.Tx{[]byte("tx1")},
			blockHeight:   1,
			timestamp:     time.Now().UTC(),
			prevStateRoot: types.Hash{},
			expectedErr:   types.ErrEmptyStateRoot,
		},
		{
			name:          "future timestamp",
			txs:           []types.Tx{[]byte("tx1")},
			blockHeight:   1,
			timestamp:     time.Now().Add(24 * time.Hour),
			prevStateRoot: types.Hash{1, 2, 3},
			expectedErr:   types.ErrFutureBlockTime,
		},
		{
			name:          "empty transaction",
			txs:           []types.Tx{[]byte("")},
			blockHeight:   1,
			timestamp:     time.Now().UTC(),
			prevStateRoot: types.Hash{1, 2, 3},
			expectedErr:   types.ErrEmptyTx,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateRoot, maxBytes, err := s.Exec.ExecuteTxs(context.Background(), tt.txs, tt.blockHeight, tt.timestamp, tt.prevStateRoot)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
			require.NotEqual(t, types.Hash{}, stateRoot)
			require.Greater(t, maxBytes, uint64(0))
		})
	}
}

func (s *DummyTestSuite) TestInitChain() {
	t := s.T()
	tests := []struct {
		name          string
		genesisTime   time.Time
		initialHeight uint64
		chainID       string
		expectedErr   error
	}{
		{
			name:          "valid case",
			genesisTime:   time.Now().UTC(),
			initialHeight: 1,
			chainID:       "test-chain",
			expectedErr:   nil,
		},
		{
			name:          "very large initial height",
			genesisTime:   time.Now().UTC(),
			initialHeight: 1000000,
			chainID:       "test-chain",
			expectedErr:   nil,
		},
		{
			name:          "zero height",
			genesisTime:   time.Now().UTC(),
			initialHeight: 0,
			chainID:       "test-chain",
			expectedErr:   types.ErrZeroInitialHeight,
		},
		{
			name:          "empty chain ID",
			genesisTime:   time.Now().UTC(),
			initialHeight: 1,
			chainID:       "",
			expectedErr:   types.ErrEmptyChainID,
		},
		{
			name:          "future genesis time",
			genesisTime:   time.Now().Add(1 * time.Hour),
			initialHeight: 1,
			chainID:       "test-chain",
			expectedErr:   types.ErrFutureGenesisTime,
		},
		{
			name:          "invalid chain ID characters",
			genesisTime:   time.Now().UTC(),
			initialHeight: 1,
			chainID:       "@invalid",
			expectedErr:   types.ErrInvalidChainID,
		},
		{
			name:          "invalid chain ID length",
			genesisTime:   time.Now().UTC(),
			initialHeight: 1,
			chainID:       strings.Repeat("a", 50),
			expectedErr:   types.ErrChainIDTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateRoot, maxBytes, err := s.Exec.InitChain(context.Background(), tt.genesisTime, tt.initialHeight, tt.chainID)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
			require.NotEqual(t, types.Hash{}, stateRoot)
			require.Greater(t, maxBytes, uint64(0))
		})
	}
}

func (s *DummyTestSuite) TestGetTxsWithConcurrency() {
	t := s.T()
	const numGoroutines = 10
	const txsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Inject transactions concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < txsPerGoroutine; j++ {
				tx := types.Tx([]byte(fmt.Sprintf("tx-%d-%d", id, j)))
				s.TxInjector.InjectTx(tx)
			}
		}(i)
	}
	wg.Wait()

	// Verify all transactions are retrievable
	txs, err := s.Exec.GetTxs(context.Background())
	require.NoError(t, err)
	require.Len(t, txs, numGoroutines*txsPerGoroutine)

	// Verify transaction uniqueness
	txMap := make(map[string]struct{})
	for _, tx := range txs {
		txMap[string(tx)] = struct{}{}
	}
	require.Len(t, txMap, numGoroutines*txsPerGoroutine)
}

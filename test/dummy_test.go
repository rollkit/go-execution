package test

import (
	"context"
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

func TestTxRemoval(t *testing.T) {
	exec := NewDummyExecutor()

	// Generate random transactions using GetRandomTxs
	randomTxs := exec.GetRandomTxs(2)

	// Inject the random transactions into the executor
	err := exec.InjectTxs(randomTxs)
	require.NoError(t, err)

	tx1 := randomTxs[0]
	tx2 := randomTxs[1]

	// Retrieve transactions and verify
	txs, err := exec.GetTxs(context.Background())
	require.NoError(t, err)
	require.Len(t, txs, 2)
	require.Contains(t, txs, randomTxs[0])
	require.Contains(t, txs, randomTxs[1])
	// ExecuteTxs was not called, so 2 txs should still be returned
	txs, err = exec.GetTxs(context.Background())
	require.NoError(t, err)
	require.Len(t, txs, 2)
	require.Contains(t, txs, tx1)
	require.Contains(t, txs, tx2)

	state, _, err := exec.ExecuteTxs(context.Background(), []types.Tx{tx1}, 1, time.Now(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, state)

	// ExecuteTxs was called, 1 tx remaining in mempool
	txs, err = exec.GetTxs(context.Background())
	require.NoError(t, err)
	require.Len(t, txs, 1)
	require.NotContains(t, txs, tx1)
	require.Contains(t, txs, tx2)
}

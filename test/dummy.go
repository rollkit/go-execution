package test

import (
	"context"
	"time"

	"github.com/rollkit/go-execution/types"
)

// DummyExecutor is a dummy implementation of the DummyExecutor interface for testing
type DummyExecutor struct {
	stateRoot types.Hash
	maxBytes  uint64
	txs       []types.Tx
}

// NewDummyExecutor creates a new dummy DummyExecutor instance
func NewDummyExecutor() *DummyExecutor {
	return &DummyExecutor{
		stateRoot: types.Hash{1, 2, 3},
		maxBytes:  1000000,
		txs:       make([]types.Tx, 0),
	}
}

// InitChain initializes the chain state with the given genesis time, initial height, and chain ID.
// It returns the state root hash, the maximum byte size, and an error if the initialization fails.
func (e *DummyExecutor) InitChain(ctx context.Context, genesisTime time.Time, initialHeight uint64, chainID string) (types.Hash, uint64, error) {
	return e.stateRoot, e.maxBytes, nil
}

// GetTxs returns the list of transactions (types.Tx) within the DummyExecutor instance and an error if any.
func (e *DummyExecutor) GetTxs(context.Context) ([]types.Tx, error) {
	return e.txs, nil
}

// ExecuteTxs simulate execution of transactions.
func (e *DummyExecutor) ExecuteTxs(ctx context.Context, txs []types.Tx, blockHeight uint64, timestamp time.Time, prevStateRoot types.Hash) (types.Hash, uint64, error) {
	e.txs = append(e.txs, txs...)
	return e.stateRoot, e.maxBytes, nil
}

// SetFinal marks block at given height as finalized. Currently not implemented.
func (e *DummyExecutor) SetFinal(ctx context.Context, blockHeight uint64) error {
	return nil
}

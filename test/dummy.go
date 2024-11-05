package test

import (
	"time"

	"github.com/rollkit/go-execution/types"
)

// Execute is a dummy implementation of the Execute interface for testing
type Execute struct {
	stateRoot types.Hash
	maxBytes  uint64
	txs       []types.Tx
}

// NewExecute creates a new dummy Execute instance
func NewExecute() *Execute {
	return &Execute{
		stateRoot: types.Hash{1, 2, 3},
		maxBytes:  1000000,
		txs:       make([]types.Tx, 0),
	}
}

// InitChain initializes the chain state with the given genesis time, initial height, and chain ID.
// It returns the state root hash, the maximum byte size, and an error if the initialization fails.
func (e *Execute) InitChain(genesisTime time.Time, initialHeight uint64, chainID string) (types.Hash, uint64, error) {
	return e.stateRoot, e.maxBytes, nil
}

// GetTxs returns the list of transactions (types.Tx) within the Execute instance and an error if any.
func (e *Execute) GetTxs() ([]types.Tx, error) {
	return e.txs, nil
}

// ExecuteTxs simulate execution of transactions.
func (e *Execute) ExecuteTxs(txs []types.Tx, blockHeight uint64, timestamp time.Time, prevStateRoot types.Hash) (types.Hash, uint64, error) {
	e.txs = append(e.txs, txs...)
	return e.stateRoot, e.maxBytes, nil
}

// SetFinal marks block at given height as finalized. Currently not implemented.
func (e *Execute) SetFinal(blockHeight uint64) error {
	return nil
}

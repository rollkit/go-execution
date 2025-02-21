package execution

import (
	"context"
	"time"

	"github.com/rollkit/go-execution/types"
)

// Executor defines the interface that execution clients must implement to be compatible with Rollkit.
// This interface enables the separation between consensus and execution layers, allowing for modular
// and pluggable execution environments.
type Executor interface {
	// InitChain initializes a new blockchain with the given genesis parameters.
	// It returns the initial state root hash and maximum allowed transaction bytes.
	// - genesisTime: The official starting time of the blockchain
	// - initialHeight: The starting block height
	// - chainID: Unique identifier for the blockchain
	InitChain(ctx context.Context, genesisTime time.Time, initialHeight uint64, chainID string) (stateRoot types.Hash, maxBytes uint64, err error)

	// GetTxs retrieves pending transactions from the execution client's mempool.
	// These transactions are candidates for inclusion in the next block.
	GetTxs(ctx context.Context) ([]types.Tx, error)

	// ExecuteTxs processes a batch of transactions to create a new block.
	// It applies the transactions sequentially and returns the new state root.
	// - txs: List of transactions to execute
	// - blockHeight: Height of the block being created
	// - timestamp: Block timestamp
	// - prevStateRoot: State root from the previous block
	ExecuteTxs(ctx context.Context, txs []types.Tx, blockHeight uint64, timestamp time.Time, prevStateRoot types.Hash) (updatedStateRoot types.Hash, maxBytes uint64, err error)

	// SetFinal marks a block as finalized at the specified height.
	// This indicates the block can no longer be reverted.
	SetFinal(ctx context.Context, blockHeight uint64) error
}

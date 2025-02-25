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
	// InitChain initializes a new blockchain instance with genesis parameters.
	// Requirements:
	// - Must generate initial state root representing empty/genesis state
	// - Must validate and store genesis parameters for future reference
	// - Must ensure idempotency (repeated calls with identical parameters should return same results)
	// - Must return error if genesis parameters are invalid
	// - Must return maxBytes indicating maximum allowed bytes for a set of transactions in a block
	//
	// Parameters:
	// - ctx: Context for timeout/cancellation control
	// - genesisTime: timestamp marking chain start time in UTC
	// - initialHeight: First block height (must be > 0)
	// - chainID: Unique identifier string for the blockchain
	//
	// Returns:
	// - stateRoot: Hash representing initial state
	// - maxBytes: Maximum allowed bytes for transacitons in a block
	// - err: Any initialization errors
	InitChain(ctx context.Context, genesisTime time.Time, initialHeight uint64, chainID string) (stateRoot types.Hash, maxBytes uint64, err error)

	// GetTxs fetches available transactions from the execution layer's mempool.
	// Requirements:
	// - Must return currently valid transactions only
	// - Must handle empty mempool case gracefully
	// - Must respect context cancellation/timeout
	// - Should perform basic transaction validation
	// - Should not remove transactions from mempool
	// - May remove invalid transactions from mempool
	//
	// Parameters:
	// - ctx: Context for timeout/cancellation control
	//
	// Returns:
	// - []types.Tx: Slice of valid transactions
	// - error: Any errors during transaction retrieval
	GetTxs(ctx context.Context) ([]types.Tx, error)

	// ExecuteTxs processes transactions to produce a new block state.
	// Requirements:
	// - Must validate state transition against previous state root
	// - Must handle empty transaction list
	// - Must maintain deterministic execution
	// - Must respect context cancellation/timeout
	// - The rest of the rules are defined by the specific execution layer
	//
	// Parameters:
	// - ctx: Context for timeout/cancellation control
	// - txs: Ordered list of transactions to execute
	// - blockHeight: Height of block being created (must be > 0)
	// - timestamp: Block creation time in UTC
	// - prevStateRoot: Previous block's state root hash
	//
	// Returns:
	// - updatedStateRoot: New state root after executing transactions
	// - maxBytes: Maximum allowed transaction size (may change with protocol updates)
	// - err: Any execution errors
	ExecuteTxs(ctx context.Context, txs []types.Tx, blockHeight uint64, timestamp time.Time, prevStateRoot types.Hash) (updatedStateRoot types.Hash, maxBytes uint64, err error)

	// SetFinal marks a block as finalized at the specified height.
	// Requirements:
	// - Must verify block exists at specified height
	// - Must be idempotent
	// - Must maintain finality guarantees (no reverting finalized blocks)
	// - Must respect context cancellation/timeout
	// - Should clean up any temporary state/resources
	//
	// Parameters:
	// - ctx: Context for timeout/cancellation control
	// - blockHeight: Height of block to finalize
	//
	// Returns:
	// - error: Any errors during finalization
	SetFinal(ctx context.Context, blockHeight uint64) error
}

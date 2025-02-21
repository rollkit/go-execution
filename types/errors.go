package types

import "errors"

var (
	// Chain initialization errors
	ErrZeroInitialHeight = errors.New("initial height cannot be zero")
	ErrEmptyChainID      = errors.New("chain ID cannot be empty")
	ErrInvalidChainID    = errors.New("chain ID contains invalid characters")
	ErrChainIDTooLong    = errors.New("chain ID exceeds maximum length")
	ErrFutureGenesisTime = errors.New("genesis time cannot be in the future")

	// Transaction execution errors
	ErrEmptyStateRoot     = errors.New("previous state root cannot be empty")
	ErrFutureBlockTime    = errors.New("block timestamp cannot be in the future")
	ErrInvalidBlockHeight = errors.New("invalid block height")
	ErrTxTooLarge         = errors.New("transaction size exceeds maximum allowed")
	ErrEmptyTx            = errors.New("transaction cannot be empty")

	// Block finalization errors
	ErrBlockNotFound      = errors.New("block not found")
	ErrBlockAlreadyExists = errors.New("block already exists")
	ErrNonSequentialBlock = errors.New("non-sequential block height")

	// Transaction pool errors
	ErrTxAlreadyExists = errors.New("transaction already exists in pool")
	ErrTxPoolFull      = errors.New("transaction pool is full")
	ErrInvalidTxFormat = errors.New("invalid transaction format")

	// Context errors
	ErrContextCanceled = errors.New("context canceled")
	ErrContextTimeout  = errors.New("context deadline exceeded")
)

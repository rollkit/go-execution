package types

import "errors"

var (
	// Chain initialization errors

	// ErrZeroInitialHeight is returned when the initial height is zero
	ErrZeroInitialHeight = errors.New("initial height cannot be zero")
	// ErrEmptyChainID is returned when the chain ID is empty
	ErrEmptyChainID = errors.New("chain ID cannot be empty")
	// ErrInvalidChainID is returned when the chain ID contains invalid characters
	ErrInvalidChainID = errors.New("chain ID contains invalid characters")
	// ErrChainIDTooLong is returned when the chain ID exceeds maximum length
	ErrChainIDTooLong = errors.New("chain ID exceeds maximum length")
	// ErrFutureGenesisTime is returned when the genesis time is in the future
	ErrFutureGenesisTime = errors.New("genesis time cannot be in the future")

	// Transaction execution errors

	// ErrEmptyStateRoot is returned when the previous state root is empty
	ErrEmptyStateRoot = errors.New("previous state root cannot be empty")
	// ErrFutureBlockTime is returned when the block timestamp is in the future
	ErrFutureBlockTime = errors.New("block timestamp cannot be in the future")
	// ErrInvalidBlockHeight is returned when the block height is invalid
	ErrInvalidBlockHeight = errors.New("invalid block height")
	// ErrTxTooLarge is returned when the transaction size exceeds maximum allowed
	ErrTxTooLarge = errors.New("transaction size exceeds maximum allowed")
	// ErrEmptyTx is returned when the transaction is empty
	ErrEmptyTx = errors.New("transaction cannot be empty")

	// Block finalization errors

	// ErrBlockNotFound is returned when the block is not found
	ErrBlockNotFound = errors.New("block not found")
	// ErrBlockAlreadyExists is returned when the block already exists
	ErrBlockAlreadyExists = errors.New("block already exists")
	// ErrNonSequentialBlock is returned when the block height is not sequential
	ErrNonSequentialBlock = errors.New("non-sequential block height")

	// Transaction pool errors

	// ErrTxAlreadyExists is returned when the transaction already exists in pool
	ErrTxAlreadyExists = errors.New("transaction already exists in pool")
	// ErrTxPoolFull is returned when the transaction pool is full
	ErrTxPoolFull = errors.New("transaction pool is full")
	// ErrInvalidTxFormat is returned when the transaction format is invalid
	ErrInvalidTxFormat = errors.New("invalid transaction format")

	// Context errors

	// ErrContextCanceled is returned when the context is canceled
	ErrContextCanceled = errors.New("context canceled")
	// ErrContextTimeout is returned when the context deadline is exceeded
	ErrContextTimeout = errors.New("context deadline exceeded")
)

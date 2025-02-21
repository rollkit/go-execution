package test

import (
	"bytes"
	"context"
	"crypto/sha512"
	"regexp"
	"slices"
	"sync"
	"time"

	"github.com/rollkit/go-execution/types"
)

var validChainIDRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]*`)

// DummyExecutor is a dummy implementation of the DummyExecutor interface for testing
type DummyExecutor struct {
	mu           sync.RWMutex
	stateRoot    types.Hash
	pendingRoots map[uint64]types.Hash
	maxBytes     uint64
	injectedTxs  []types.Tx
}

func NewDummyExecutor() *DummyExecutor {
	return &DummyExecutor{
		stateRoot:    types.Hash{1, 2, 3},
		pendingRoots: make(map[uint64]types.Hash),
		maxBytes:     1000000,
	}
}

func (e *DummyExecutor) InitChain(ctx context.Context, genesisTime time.Time, initialHeight uint64, chainID string) (types.Hash, uint64, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if initialHeight == 0 {
		return types.Hash{}, 0, types.ErrZeroInitialHeight
	}
	if chainID == "" {
		return types.Hash{}, 0, types.ErrEmptyChainID
	}
	if !validChainIDRegex.MatchString(chainID) {
		return types.Hash{}, 0, types.ErrInvalidChainID
	}
	if genesisTime.After(time.Now()) {
		return types.Hash{}, 0, types.ErrFutureGenesisTime
	}
	if len(chainID) > 32 {
		return types.Hash{}, 0, types.ErrChainIDTooLong
	}

	hash := sha512.New()
	hash.Write(e.stateRoot)
	e.stateRoot = hash.Sum(nil)
	return e.stateRoot, e.maxBytes, nil
}

func (e *DummyExecutor) ExecuteTxs(ctx context.Context, txs []types.Tx, blockHeight uint64, timestamp time.Time, prevStateRoot types.Hash) (types.Hash, uint64, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if bytes.Equal(prevStateRoot, types.Hash{}) {
		return types.Hash{}, 0, types.ErrEmptyStateRoot
	}

	// Don't really allow future block times, but allow up to 5 minutes in the future
	// for testing purposes.
	if timestamp.After(time.Now().Add(5 * time.Minute)) {
		return types.Hash{}, 0, types.ErrFutureBlockTime
	}
	if blockHeight == 0 {
		return types.Hash{}, 0, types.ErrInvalidBlockHeight
	}

	for _, tx := range txs {
		if len(tx) == 0 {
			return types.Hash{}, 0, types.ErrEmptyTx
		}
		if uint64(len(tx)) > e.maxBytes {
			return types.Hash{}, 0, types.ErrTxTooLarge
		}
	}

	hash := sha512.New()
	hash.Write(prevStateRoot)
	for _, tx := range txs {
		hash.Write(tx)
	}
	pending := hash.Sum(nil)
	e.pendingRoots[blockHeight] = pending
	e.removeExecutedTxs(txs)
	return pending, e.maxBytes, nil
}

func (e *DummyExecutor) SetFinal(ctx context.Context, blockHeight uint64) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if pending, ok := e.pendingRoots[blockHeight]; ok {
		e.stateRoot = pending
		delete(e.pendingRoots, blockHeight)
		return nil
	}
	return types.ErrBlockNotFound
}

func (e *DummyExecutor) GetTxs(context.Context) ([]types.Tx, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	txs := make([]types.Tx, len(e.injectedTxs))
	copy(txs, e.injectedTxs)
	return txs, nil
}

func (e *DummyExecutor) InjectTx(tx types.Tx) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.injectedTxs = append(e.injectedTxs, tx)
}

func (e *DummyExecutor) removeExecutedTxs(txs []types.Tx) {
	e.injectedTxs = slices.DeleteFunc(e.injectedTxs, func(tx types.Tx) bool {
		return slices.ContainsFunc(txs, func(t types.Tx) bool { return bytes.Equal(tx, t) })
	})
}

func (e *DummyExecutor) GetStateRoot() types.Hash {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.stateRoot
}

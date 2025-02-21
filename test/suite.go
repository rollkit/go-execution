package test

import (
	"context"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/rollkit/go-execution"
	"github.com/rollkit/go-execution/types"
)

// ExecutorSuite is a reusable test suite for Execution API implementations.
type ExecutorSuite struct {
	suite.Suite
	Exec       execution.Executor
	TxInjector TxInjector
}

// TxInjector provides an interface for injecting transactions into a test suite.
type TxInjector interface {
	InjectTx(tx types.Tx)
}

// TestInitChain tests InitChain method.
func (s *ExecutorSuite) TestInitChain() {
	genesisTime := time.Now().UTC()
	initialHeight := uint64(1)
	chainID := "test-chain"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stateRoot, maxBytes, err := s.Exec.InitChain(ctx, genesisTime, initialHeight, chainID)
	s.Require().NoError(err)
	s.NotEqual(types.Hash{}, stateRoot)
	s.Greater(maxBytes, uint64(0))
}

// TestGetTxs tests GetTxs method.
func (s *ExecutorSuite) TestGetTxs() {
	s.skipIfInjectorNotSet()

	tx1 := types.Tx("tx1")
	tx2 := types.Tx("tx2")

	s.TxInjector.InjectTx(tx1)
	s.TxInjector.InjectTx(tx2)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	txs, err := s.Exec.GetTxs(ctx)
	s.Require().NoError(err)
	s.Require().Len(txs, 2)
	s.Require().Contains(txs, tx1)
	s.Require().Contains(txs, tx2)
}

func (s *ExecutorSuite) skipIfInjectorNotSet() {
	if s.TxInjector == nil {
		s.T().Skipf("Skipping %s because TxInjector is not provided", s.T().Name())
	}
}

// TestExecuteTxs tests ExecuteTxs method.
func (s *ExecutorSuite) TestExecuteTxs() {
	txs := []types.Tx{[]byte("tx1"), []byte("tx2")}
	blockHeight := uint64(1)
	timestamp := time.Now().UTC()
	prevStateRoot := types.Hash{1, 2, 3}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stateRoot, maxBytes, err := s.Exec.ExecuteTxs(ctx, txs, blockHeight, timestamp, prevStateRoot)
	s.Require().NoError(err)
	s.NotEqual(types.Hash{}, stateRoot)
	s.Greater(maxBytes, uint64(0))
}

// TestSetFinal tests SetFinal method.
func (s *ExecutorSuite) TestSetFinal() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// finalizing invalid height must return error
	err := s.Exec.SetFinal(ctx, 1)
	s.Require().Error(err)

	ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel2()
	_, _, err = s.Exec.ExecuteTxs(ctx2, nil, 2, time.Now(), types.Hash("test state"))
	s.Require().NoError(err)

	ctx3, cancel3 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel3()
	err = s.Exec.SetFinal(ctx3, 2)
	s.Require().NoError(err)
}

// TestMultipleBlocks is a basic test ensuring that all API methods used together can be used to produce multiple blocks.
func (s *ExecutorSuite) TestMultipleBlocks() {
	genesisTime := time.Now().UTC()
	initialHeight := uint64(1)
	chainID := "test-chain"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stateRoot, maxBytes, err := s.Exec.InitChain(ctx, genesisTime, initialHeight, chainID)
	s.Require().NoError(err)
	s.NotEqual(types.Hash{}, stateRoot)
	s.Greater(maxBytes, uint64(0))

	for i := initialHeight; i <= 10; i++ {
		txs, err := s.Exec.GetTxs(ctx)
		s.Require().NoError(err)

		blockTime := genesisTime.Add(time.Duration(i+1) * time.Second) //nolint:gosec
		stateRoot, maxBytes, err = s.Exec.ExecuteTxs(ctx, txs, i, blockTime, stateRoot)
		s.Require().NoError(err)
		s.Require().NotZero(maxBytes)

		err = s.Exec.SetFinal(ctx, i)
		s.Require().NoError(err)
	}
}

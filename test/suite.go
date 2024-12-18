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
	Exec execution.Executor
}

// TestInitChain tests InitChain method.
func (s *ExecutorSuite) TestInitChain() {
	genesisTime := time.Now().UTC()
	initialHeight := uint64(1)
	chainID := "test-chain"

	stateRoot, maxBytes, err := s.Exec.InitChain(context.TODO(), genesisTime, initialHeight, chainID)
	s.Require().NoError(err)
	s.NotEqual(types.Hash{}, stateRoot)
	s.Greater(maxBytes, uint64(0))
}

// TestGetTxs tests GetTxs method.
func (s *ExecutorSuite) TestGetTxs() {
	txs, err := s.Exec.GetTxs(context.TODO())
	s.Require().NoError(err)
	s.Empty(txs)
}

// TestExecuteTxs tests ExecuteTxs method.
func (s *ExecutorSuite) TestExecuteTxs() {
	txs := []types.Tx{[]byte("tx1"), []byte("tx2")}
	blockHeight := uint64(1)
	timestamp := time.Now().UTC()
	prevStateRoot := types.Hash{1, 2, 3}

	stateRoot, maxBytes, err := s.Exec.ExecuteTxs(context.TODO(), txs, blockHeight, timestamp, prevStateRoot)
	s.Require().NoError(err)
	s.NotEqual(types.Hash{}, stateRoot)
	s.Greater(maxBytes, uint64(0))
}

// TestSetFinal tests SetFinal method.
func (s *ExecutorSuite) TestSetFinal() {
	// finalizing invalid height must return error
	err := s.Exec.SetFinal(context.TODO(), 1)
	s.Require().Error(err)

	_, _, err = s.Exec.ExecuteTxs(context.TODO(), nil, 2, time.Now(), types.Hash("test state"))
	s.Require().NoError(err)
	err = s.Exec.SetFinal(context.TODO(), 2)
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

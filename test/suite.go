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
	InjectRandomTx() types.Tx
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// try to get transactions without injecting any
	txs, err := s.Exec.GetTxs(ctx)
	s.Require().NoError(err)
	s.Require().Empty(txs)

	// inject two txs and retrieve them
	tx1 := s.TxInjector.InjectRandomTx()
	tx2 := s.TxInjector.InjectRandomTx()
	txs, err = s.Exec.GetTxs(ctx)
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
	s.skipIfInjectorNotSet()

	cases := []struct {
		name string
		txs  []types.Tx
	}{
		{
			name: "nil txs",
			txs:  nil,
		},
		{
			name: "empty txs",
			txs:  []types.Tx{},
		},
		{
			name: "two txs",
			txs:  []types.Tx{s.TxInjector.InjectRandomTx(), s.TxInjector.InjectRandomTx()},
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			initialHeight := uint64(1)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			genesisTime, genesisStateRoot, _ := s.initChain(ctx, initialHeight)

			stateRoot, maxBytes, err := s.Exec.ExecuteTxs(ctx, c.txs, initialHeight, genesisTime.Add(time.Second), genesisStateRoot)
			s.Require().NoError(err)
			s.Require().NotEmpty(stateRoot)
			s.Require().NotEqual(genesisStateRoot, stateRoot)
			s.Require().Greater(maxBytes, uint64(0))
		})
	}
}

// TestSetFinal tests SetFinal method.
func (s *ExecutorSuite) TestSetFinal() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// finalizing invalid height must return error
	err := s.Exec.SetFinal(ctx, 7)
	s.Require().Error(err)

	initialHeight := uint64(1)
	_, stateRoot, _ := s.initChain(ctx, initialHeight)
	_, _, err = s.Exec.ExecuteTxs(ctx, nil, initialHeight, time.Now(), stateRoot)
	s.Require().NoError(err)
	err = s.Exec.SetFinal(ctx, initialHeight)
	s.Require().NoError(err)
}

// TestMultipleBlocks is a basic test ensuring that all API methods used together can be used to produce multiple blocks.
func (s *ExecutorSuite) TestMultipleBlocks() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	initialHeight := uint64(1)
	genesisTime, prevStateRoot, _ := s.initChain(ctx, initialHeight)

	for i := initialHeight; i <= 10; i++ {
		txs, err := s.Exec.GetTxs(ctx)
		s.Require().NoError(err)

		blockTime := genesisTime.Add(time.Duration(i+1) * time.Second) //nolint:gosec
		stateRoot, maxBytes, err := s.Exec.ExecuteTxs(ctx, txs, i, blockTime, prevStateRoot)
		s.Require().NoError(err)
		s.Require().NotZero(maxBytes)
		s.Require().NotEqual(prevStateRoot, stateRoot)

		prevStateRoot = stateRoot

		err = s.Exec.SetFinal(ctx, i)
		s.Require().NoError(err)
	}
}

func (s *ExecutorSuite) initChain(ctx context.Context, initialHeight uint64) (time.Time, types.Hash, uint64) {
	genesisTime := time.Now().UTC()
	chainID := "test-chain"

	stateRoot, maxBytes, err := s.Exec.InitChain(ctx, genesisTime, initialHeight, chainID)
	s.Require().NoError(err)
	s.Require().NotEmpty(stateRoot)
	return genesisTime, stateRoot, maxBytes
}

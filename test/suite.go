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

	stateRoot, maxBytes, err := s.Exec.InitChain(context.TODO(), genesisTime, initialHeight, chainID)
	s.Require().NoError(err)
	s.NotEqual(types.Hash{}, stateRoot)
	s.Greater(maxBytes, uint64(0))
}

// TestGetTxs tests GetTxs method.
func (s *ExecutorSuite) TestGetTxs() {
	s.skipIfInjectorNotSet()

	tx1 := s.TxInjector.InjectRandomTx()
	tx2 := s.TxInjector.InjectRandomTx()
	txs, err := s.Exec.GetTxs(context.TODO())
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
	txs := []types.Tx{s.TxInjector.InjectRandomTx(), s.TxInjector.InjectRandomTx()}

	genesisTime, initialHeight, genesisStateRoot, _ := s.initChain(context.TODO())

	stateRoot, maxBytes, err := s.Exec.ExecuteTxs(context.TODO(), txs, initialHeight+1, genesisTime.Add(time.Second), genesisStateRoot)
	s.Require().NoError(err)
	s.NotEqual(types.Hash{}, stateRoot)
	s.NotEqual(genesisStateRoot, stateRoot)
	s.Greater(maxBytes, uint64(0))
}

// TestSetFinal tests SetFinal method.
func (s *ExecutorSuite) TestSetFinal() {
	// finalizing invalid height must return error
	err := s.Exec.SetFinal(context.TODO(), 1)
	s.Require().Error(err)

	_, _, _, _ = s.initChain(context.TODO())
	_, _, err = s.Exec.ExecuteTxs(context.TODO(), nil, 2, time.Now(), types.Hash("test state"))
	s.Require().NoError(err)
	err = s.Exec.SetFinal(context.TODO(), 2)
	s.Require().NoError(err)
}

// TestMultipleBlocks is a basic test ensuring that all API methods used together can be used to produce multiple blocks.
func (s *ExecutorSuite) TestMultipleBlocks() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	genesisTime, initialHeight, stateRoot, maxBytes := s.initChain(ctx)

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

func (s *ExecutorSuite) initChain(ctx context.Context) (time.Time, uint64, types.Hash, uint64) {
	genesisTime := time.Now().UTC()
	initialHeight := uint64(1)
	chainID := "test-chain"

	stateRoot, maxBytes, err := s.Exec.InitChain(ctx, genesisTime, initialHeight, chainID)
	s.Require().NoError(err)
	return genesisTime, initialHeight, stateRoot, maxBytes
}

package test

import (
	"testing"
	"time"

	"github.com/rollkit/go-execution/types"
	"github.com/stretchr/testify/suite"

	"github.com/rollkit/go-execution"
)

type ExecuteSuite struct {
	suite.Suite
	Exec execution.Execute
}

func (s *ExecuteSuite) SetupTest() {}

func (s *ExecuteSuite) TestInitChain() {
	genesisTime := time.Now().UTC()
	initialHeight := uint64(1)
	chainID := "test-chain"

	stateRoot, maxBytes, err := s.Exec.InitChain(genesisTime, initialHeight, chainID)
	s.Require().NoError(err)
	s.NotEqual(types.Hash{}, stateRoot)
	s.Greater(maxBytes, uint64(0))
}

func (s *ExecuteSuite) TestGetTxs() {
	txs, err := s.Exec.GetTxs()
	s.Require().NoError(err)
	s.NotNil(txs)
}

func (s *ExecuteSuite) TestExecuteTxs() {
	txs := []types.Tx{[]byte("tx1"), []byte("tx2")}
	blockHeight := uint64(1)
	timestamp := time.Now().UTC()
	prevStateRoot := types.Hash{1, 2, 3}

	stateRoot, maxBytes, err := s.Exec.ExecuteTxs(txs, blockHeight, timestamp, prevStateRoot)
	s.Require().NoError(err)
	s.NotEqual(types.Hash{}, stateRoot)
	s.Greater(maxBytes, uint64(0))
}

func (s *ExecuteSuite) TestSetFinal() {
	err := s.Exec.SetFinal(1)
	s.Require().NoError(err)
}

type DummyTestSuite struct {
	ExecuteSuite
}

func (s *DummyTestSuite) SetupTest() {
	s.Exec = NewExecute()
}

func TestDummySuite(t *testing.T) {
	suite.Run(t, new(DummyTestSuite))
}

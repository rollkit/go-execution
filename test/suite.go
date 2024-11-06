package test

import (
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/rollkit/go-execution"
	"github.com/rollkit/go-execution/types"
)

// ExecuteSuite is a reusable test suite for Execution API implementations.
type ExecuteSuite struct {
	suite.Suite
	Exec execution.Execute
}

// TestInitChain tests InitChain method.
func (s *ExecuteSuite) TestInitChain() {
	genesisTime := time.Now().UTC()
	initialHeight := uint64(1)
	chainID := "test-chain"

	stateRoot, maxBytes, err := s.Exec.InitChain(genesisTime, initialHeight, chainID)
	s.Require().NoError(err)
	s.NotEqual(types.Hash{}, stateRoot)
	s.Greater(maxBytes, uint64(0))
}

// TestGetTxs tests GetTxs method.
func (s *ExecuteSuite) TestGetTxs() {
	txs, err := s.Exec.GetTxs()
	s.Require().NoError(err)
	s.NotNil(txs)
}

// TestExecuteTxs tests ExecuteTxs method.
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

// TestSetFinal tests SetFinal method.
func (s *ExecuteSuite) TestSetFinal() {
	err := s.Exec.SetFinal(1)
	s.Require().NoError(err)
}

package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DummyTestSuite struct {
	ExecutorSuite
}

func (s *DummyTestSuite) SetupTest() {
	s.Exec = NewDummyExecutor()
}

func TestDummySuite(t *testing.T) {
	suite.Run(t, new(DummyTestSuite))
}

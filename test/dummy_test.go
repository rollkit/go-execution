package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DummyTestSuite struct {
	ExecuteSuite
}

func (s *DummyTestSuite) SetupTest() {
	s.Exec = NewExecute()
}

func TestDummySuite(t *testing.T) {
	suite.Run(t, new(DummyTestSuite))
}

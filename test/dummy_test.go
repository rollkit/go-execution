package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestDummySuite(t *testing.T) {
	suite.Run(t, new(DummyTestSuite))
}

package types

import (
	"github.com/celestiaorg/go-header"
)

// Tx represents a transaction in the form of a byte slice.
type Tx []byte

// Txs represents a slice of transactions.
type Txs []Tx

// Hash is a type alias for header.Hash
type Hash = header.Hash

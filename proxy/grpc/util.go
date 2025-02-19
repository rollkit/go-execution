package grpc

import "github.com/rollkit/go-execution/types"

func copyHash(src []byte) types.Hash {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

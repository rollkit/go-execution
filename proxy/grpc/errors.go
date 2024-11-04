package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnknownPayload      = status.Error(codes.NotFound, "payload does not exist")
	ErrInvalidForkchoice   = status.Error(codes.InvalidArgument, "invalid forkchoice state")
	ErrInvalidPayloadAttrs = status.Error(codes.InvalidArgument, "invalid payload attributes")
	ErrTooLargeRequest     = status.Error(codes.ResourceExhausted, "request too large")
	ErrUnsupportedFork     = status.Error(codes.Unimplemented, "unsupported fork")
	ErrInvalidJWT          = status.Error(codes.Unauthenticated, "invalid JWT token")
)

package jsonrpc

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	ErrCodeParse          = -32700
	ErrCodeInvalidRequest = -32600
	ErrCodeMethodNotFound = -32601
	ErrCodeInvalidParams  = -32602
	ErrCodeInternal       = -32603
)

var (
	ErrParse          = &jsonRPCError{Code: ErrCodeParse, Message: "Parse error"}
	ErrInvalidRequest = &jsonRPCError{Code: ErrCodeInvalidRequest, Message: "Invalid request"}
	ErrMethodNotFound = &jsonRPCError{Code: ErrCodeMethodNotFound, Message: "Method not found"}
	ErrInvalidParams  = &jsonRPCError{Code: ErrCodeInvalidParams, Message: "Invalid params"}
	ErrInternal       = &jsonRPCError{Code: ErrCodeInternal, Message: "Internal error"}
)

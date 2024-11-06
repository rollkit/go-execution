package jsonrpc

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	// ErrCodeParse is a reserved JSON-RPC error code
	ErrCodeParse = -32700
	// ErrCodeInvalidRequest is a reserved JSON-RPC error code
	ErrCodeInvalidRequest = -32600
	// ErrCodeMethodNotFound is a reserved JSON-RPC error code
	ErrCodeMethodNotFound = -32601
	// ErrCodeInvalidParams is a reserved JSON-RPC error code
	ErrCodeInvalidParams = -32602
	// ErrCodeInternal is a reserved JSON-RPC error code
	ErrCodeInternal = -32603
)

var (

	// ErrParse represents a JSON-RPC error indicating a problem with parsing the JSON request payload.
	ErrParse = &jsonRPCError{Code: ErrCodeParse, Message: "Parse error"}

	// ErrInvalidRequest represents a JSON-RPC error indicating an invalid JSON request payload.
	ErrInvalidRequest = &jsonRPCError{Code: ErrCodeInvalidRequest, Message: "Invalid request"}

	// ErrMethodNotFound represents a JSON-RPC error indicating that the requested method could not be found.
	ErrMethodNotFound = &jsonRPCError{Code: ErrCodeMethodNotFound, Message: "Method not found"}

	// ErrInvalidParams represents a JSON-RPC error indicating invalid parameters in the request.
	ErrInvalidParams = &jsonRPCError{Code: ErrCodeInvalidParams, Message: "Invalid params"}

	// ErrInternal represents a JSON-RPC error indicating an unspecified internal error within the server.
	ErrInternal = &jsonRPCError{Code: ErrCodeInternal, Message: "Internal error"}
)

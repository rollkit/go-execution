package jsonrpc

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rollkit/go-execution"
	"github.com/rollkit/go-execution/types"
)

type Server struct {
	exec   execution.Execute
	config *Config
}

func NewServer(exec execution.Execute, config *Config) *Server {
	if config == nil {
		config = DefaultConfig()
	}
	return &Server{
		exec:   exec,
		config: config,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, ErrInvalidRequest)
		return
	}

	if r.ContentLength > s.config.MaxRequestSize {
		writeError(w, ErrInvalidRequest)
		return
	}

	var request struct {
		JSONRPC string          `json:"jsonrpc"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params"`
		ID      interface{}     `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, ErrParse)
		return
	}

	if request.JSONRPC != "2.0" {
		writeError(w, ErrInvalidRequest)
		return
	}

	var result interface{}
	var err *jsonRPCError

	switch request.Method {
	case "init_chain":
		result, err = s.handleInitChain(request.Params)
	case "get_txs":
		result, err = s.handleGetTxs()
	case "execute_txs":
		result, err = s.handleExecuteTxs(request.Params)
	case "set_final":
		result, err = s.handleSetFinal(request.Params)
	default:
		err = ErrMethodNotFound
	}

	if err != nil {
		writeResponse(w, request.ID, nil, err)
		return
	}

	writeResponse(w, request.ID, result, nil)
}

func (s *Server) handleInitChain(params json.RawMessage) (interface{}, *jsonRPCError) {
	var p struct {
		GenesisTime   int64  `json:"genesis_time"`
		InitialHeight uint64 `json:"initial_height"`
		ChainID       string `json:"chain_id"`
	}

	if err := json.Unmarshal(params, &p); err != nil {
		return nil, ErrInvalidParams
	}

	stateRoot, maxBytes, err := s.exec.InitChain(
		time.Unix(p.GenesisTime, 0).UTC(),
		p.InitialHeight,
		p.ChainID,
	)
	if err != nil {
		return nil, &jsonRPCError{Code: ErrCodeInternal, Message: err.Error()}
	}

	return map[string]interface{}{
		"state_root": base64.StdEncoding.EncodeToString(stateRoot[:]),
		"max_bytes":  maxBytes,
	}, nil
}

func (s *Server) handleGetTxs() (interface{}, *jsonRPCError) {
	txs, err := s.exec.GetTxs()
	if err != nil {
		return nil, &jsonRPCError{Code: ErrCodeInternal, Message: err.Error()}
	}

	encodedTxs := make([]string, len(txs))
	for i, tx := range txs {
		encodedTxs[i] = base64.StdEncoding.EncodeToString(tx)
	}

	return map[string]interface{}{
		"txs": encodedTxs,
	}, nil
}

func (s *Server) handleExecuteTxs(params json.RawMessage) (interface{}, *jsonRPCError) {
	var p struct {
		Txs           []string `json:"txs"`
		BlockHeight   uint64   `json:"block_height"`
		Timestamp     int64    `json:"timestamp"`
		PrevStateRoot string   `json:"prev_state_root"`
	}

	if err := json.Unmarshal(params, &p); err != nil {
		return nil, ErrInvalidParams
	}

	// Decode base64 txs
	txs := make([]types.Tx, len(p.Txs))
	for i, encodedTx := range p.Txs {
		tx, err := base64.StdEncoding.DecodeString(encodedTx)
		if err != nil {
			return nil, ErrInvalidParams
		}
		txs[i] = tx
	}

	// Decode base64 prev state root
	prevStateRootBytes, err := base64.StdEncoding.DecodeString(p.PrevStateRoot)
	if err != nil {
		return nil, ErrInvalidParams
	}

	var prevStateRoot types.Hash
	copy(prevStateRoot[:], prevStateRootBytes)

	updatedStateRoot, maxBytes, err := s.exec.ExecuteTxs(
		txs,
		p.BlockHeight,
		time.Unix(p.Timestamp, 0).UTC(),
		prevStateRoot,
	)
	if err != nil {
		return nil, &jsonRPCError{Code: ErrCodeInternal, Message: err.Error()}
	}

	return map[string]interface{}{
		"updated_state_root": base64.StdEncoding.EncodeToString(updatedStateRoot[:]),
		"max_bytes":          maxBytes,
	}, nil
}

func (s *Server) handleSetFinal(params json.RawMessage) (interface{}, *jsonRPCError) {
	var p struct {
		BlockHeight uint64 `json:"block_height"`
	}

	if err := json.Unmarshal(params, &p); err != nil {
		return nil, ErrInvalidParams
	}

	if err := s.exec.SetFinal(p.BlockHeight); err != nil {
		return nil, &jsonRPCError{Code: ErrCodeInternal, Message: err.Error()}
	}

	return map[string]interface{}{}, nil
}

func writeError(w http.ResponseWriter, err *jsonRPCError) {
	writeResponse(w, nil, nil, err)
}

func writeResponse(w http.ResponseWriter, id interface{}, result interface{}, err *jsonRPCError) {
	response := struct {
		JSONRPC string        `json:"jsonrpc"`
		Result  interface{}   `json:"result,omitempty"`
		Error   *jsonRPCError `json:"error,omitempty"`
		ID      interface{}   `json:"id"`
	}{
		JSONRPC: "2.0",
		Result:  result,
		Error:   err,
		ID:      id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

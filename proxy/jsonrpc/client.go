package jsonrpc

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rollkit/go-execution/types"
)

// Client defines JSON-RPC proxy client of execution API.
type Client struct {
	endpoint string
	client   *http.Client
	config   *Config
}

// NewClient creates new proxy client with default config.
func NewClient() *Client {
	return &Client{
		config: DefaultConfig(),
		client: &http.Client{},
	}
}

// SetConfig updates the client's configuration with the provided config.
func (c *Client) SetConfig(config *Config) {
	if config != nil {
		c.config = config
		c.client.Timeout = config.DefaultTimeout
	}
}

// Start is used to start the client.
func (c *Client) Start(endpoint string) error {
	c.endpoint = endpoint
	return nil
}

// Stop method is used to stop the client.
func (c *Client) Stop() error {
	return nil
}

// InitChain initializes the blockchain with genesis information.
func (c *Client) InitChain(ctx context.Context, genesisTime time.Time, initialHeight uint64, chainID string) (types.Hash, uint64, error) {
	params := map[string]interface{}{
		"genesis_time":   genesisTime.Unix(),
		"initial_height": initialHeight,
		"chain_id":       chainID,
	}

	var result struct {
		StateRoot string `json:"state_root"`
		MaxBytes  uint64 `json:"max_bytes"`
	}

	if err := c.call(context.TODO(), "init_chain", params, &result); err != nil {
		return types.Hash{}, 0, err
	}

	stateRootBytes, err := base64.StdEncoding.DecodeString(result.StateRoot)
	if err != nil {
		return types.Hash{}, 0, fmt.Errorf("failed to decode state root: %w", err)
	}

	var stateRoot types.Hash
	copy(stateRoot[:], stateRootBytes)

	return stateRoot, result.MaxBytes, nil
}

// GetTxs retrieves all available transactions from the execution client's mempool.
func (c *Client) GetTxs(context.Context) ([]types.Tx, error) {
	var result struct {
		Txs []string `json:"txs"`
	}

	if err := c.call(context.TODO(), "get_txs", nil, &result); err != nil {
		return nil, err
	}

	txs := make([]types.Tx, len(result.Txs))
	for i, encodedTx := range result.Txs {
		tx, err := base64.StdEncoding.DecodeString(encodedTx)
		if err != nil {
			return nil, fmt.Errorf("failed to decode tx: %w", err)
		}
		txs[i] = tx
	}

	return txs, nil
}

// ExecuteTxs executes a set of transactions to produce a new block header.
func (c *Client) ExecuteTxs(ctx context.Context, txs []types.Tx, blockHeight uint64, timestamp time.Time, prevStateRoot types.Hash) (types.Hash, uint64, error) {
	// Encode txs to base64
	encodedTxs := make([]string, len(txs))
	for i, tx := range txs {
		encodedTxs[i] = base64.StdEncoding.EncodeToString(tx)
	}

	params := map[string]interface{}{
		"txs":             encodedTxs,
		"block_height":    blockHeight,
		"timestamp":       timestamp.Unix(),
		"prev_state_root": base64.StdEncoding.EncodeToString(prevStateRoot[:]),
	}

	var result struct {
		UpdatedStateRoot string `json:"updated_state_root"`
		MaxBytes         uint64 `json:"max_bytes"`
	}

	if err := c.call(context.TODO(), "execute_txs", params, &result); err != nil {
		return types.Hash{}, 0, err
	}

	updatedStateRootBytes, err := base64.StdEncoding.DecodeString(result.UpdatedStateRoot)
	if err != nil {
		return types.Hash{}, 0, fmt.Errorf("failed to decode updated state root: %w", err)
	}

	var updatedStateRoot types.Hash
	copy(updatedStateRoot[:], updatedStateRootBytes)

	return updatedStateRoot, result.MaxBytes, nil
}

// SetFinal marks a block at the given height as final.
func (c *Client) SetFinal(ctx context.Context, blockHeight uint64) error {
	params := map[string]interface{}{
		"block_height": blockHeight,
	}

	return c.call(context.TODO(), "set_final", params, nil)
}

func (c *Client) call(ctx context.Context, method string, params interface{}, result interface{}) error {
	request := struct {
		JSONRPC string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params,omitempty"`
		ID      int         `json:"id"`
	}{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var jsonRPCResponse struct {
		Error  *jsonRPCError   `json:"error,omitempty"`
		Result json.RawMessage `json:"result,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jsonRPCResponse); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if jsonRPCResponse.Error != nil {
		return fmt.Errorf("RPC error: %d %s", jsonRPCResponse.Error.Code, jsonRPCResponse.Error.Message)
	}

	if result != nil {
		if err := json.Unmarshal(jsonRPCResponse.Result, result); err != nil {
			return fmt.Errorf("failed to unmarshal result: %w", err)
		}
	}

	return nil
}

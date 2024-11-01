package jsonrpc

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rollkit/rollkit/types"
)

type Client struct {
	endpoint string
	client   *http.Client
	config   *Config
}

func NewClient() *Client {
	return &Client{
		config: DefaultConfig(),
		client: &http.Client{},
	}
}

func (c *Client) SetConfig(config *Config) {
	if config != nil {
		c.config = config
		c.client.Timeout = config.DefaultTimeout
	}
}

func (c *Client) Start(endpoint string) error {
	c.endpoint = endpoint
	return nil
}

func (c *Client) Stop() error {
	return nil
}

func (c *Client) InitChain(genesisTime time.Time, initialHeight uint64, chainID string) (types.Hash, uint64, error) {
	params := map[string]interface{}{
		"genesis_time":   genesisTime.Unix(),
		"initial_height": initialHeight,
		"chain_id":       chainID,
	}

	var result struct {
		StateRoot string `json:"state_root"`
		MaxBytes  uint64 `json:"max_bytes"`
	}

	if err := c.call("init_chain", params, &result); err != nil {
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

func (c *Client) GetTxs() ([]types.Tx, error) {
	var result struct {
		Txs []string `json:"txs"`
	}

	if err := c.call("get_txs", nil, &result); err != nil {
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

func (c *Client) ExecuteTxs(txs []types.Tx, blockHeight uint64, timestamp time.Time, prevStateRoot types.Hash) (types.Hash, uint64, error) {
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

	if err := c.call("execute_txs", params, &result); err != nil {
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

func (c *Client) SetFinal(blockHeight uint64) error {
	params := map[string]interface{}{
		"block_height": blockHeight,
	}

	return c.call("set_final", params, nil)
}

func (c *Client) call(method string, params interface{}, result interface{}) error {
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

	req, err := http.NewRequestWithContext(context.Background(), "POST", c.endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

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
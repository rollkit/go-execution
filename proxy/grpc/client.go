package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/rollkit/go-execution/types"
	pb "github.com/rollkit/go-execution/types/pb/execution"
)

// Client defines gRPC proxy client
type Client struct {
	conn   *grpc.ClientConn
	client pb.ExecutionServiceClient
	config *Config
}

// NewClient creates a new instance of Client with default configuration.
func NewClient() *Client {
	return &Client{
		config: DefaultConfig(),
	}
}

// SetConfig sets the configuration for the Client instance.
func (c *Client) SetConfig(config *Config) {
	if config != nil {
		c.config = config
	}
}

// Start initializes the Client by creating a new gRPC connection and storing the ExecutionServiceClient instance.
func (c *Client) Start(target string, opts ...grpc.DialOption) error {
	var err error
	c.conn, err = grpc.NewClient(target, opts...)
	if err != nil {
		return err
	}
	c.client = pb.NewExecutionServiceClient(c.conn)
	return nil
}

// Stop stops the client by closing the underlying gRPC connection if it exists.
func (c *Client) Stop() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// InitChain initializes the blockchain with genesis information.
func (c *Client) InitChain(ctx context.Context, genesisTime time.Time, initialHeight uint64, chainID string) (types.Hash, uint64, error) {
	resp, err := c.client.InitChain(ctx, &pb.InitChainRequest{
		GenesisTime:   genesisTime.Unix(),
		InitialHeight: initialHeight,
		ChainId:       chainID,
	})
	if err != nil {
		return types.Hash{}, 0, err
	}

	stateRoot := make([]byte, len(resp.StateRoot))
	copy(stateRoot[:], resp.StateRoot)

	return stateRoot, resp.MaxBytes, nil
}

// GetTxs retrieves all available transactions from the execution client's mempool.
func (c *Client) GetTxs(ctx context.Context) ([]types.Tx, error) {
	resp, err := c.client.GetTxs(ctx, &pb.GetTxsRequest{})
	if err != nil {
		return nil, err
	}

	txs := make([]types.Tx, len(resp.Txs))
	for i, tx := range resp.Txs {
		txs[i] = tx
	}

	return txs, nil
}

// ExecuteTxs executes a set of transactions to produce a new block header.
func (c *Client) ExecuteTxs(ctx context.Context, txs []types.Tx, blockHeight uint64, timestamp time.Time, prevStateRoot types.Hash) (types.Hash, uint64, error) {
	req := &pb.ExecuteTxsRequest{
		Txs:           make([][]byte, len(txs)),
		BlockHeight:   blockHeight,
		Timestamp:     timestamp.Unix(),
		PrevStateRoot: prevStateRoot[:],
	}
	for i, tx := range txs {
		req.Txs[i] = tx
	}

	resp, err := c.client.ExecuteTxs(ctx, req)
	if err != nil {
		return types.Hash{}, 0, err
	}

	updatedStateRoot := make([]byte, len(resp.UpdatedStateRoot))
	copy(updatedStateRoot[:], resp.UpdatedStateRoot)

	return updatedStateRoot, resp.MaxBytes, nil
}

// SetFinal marks a block at the given height as final.
func (c *Client) SetFinal(ctx context.Context, blockHeight uint64) error {
	_, err := c.client.SetFinal(ctx, &pb.SetFinalRequest{
		BlockHeight: blockHeight,
	})
	return err
}

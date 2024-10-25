package grpc

import (
	"context"
	"time"

	"github.com/rollkit/rollkit/types"
	"google.golang.org/grpc"

	"github.com/LastL2/execution"
	pb "github.com/LastL2/types"
)

// Client implements the execution.Execute interface using gRPC
type Client struct {
	client pb.ExecutionServiceClient
}

// NewClient creates a new gRPC client
func NewClient(conn grpc.ClientConnInterface) execution.Execute {
	return &Client{
		client: pb.NewExecutionServiceClient(conn),
	}
}

func (c *Client) InitChain(genesisTime time.Time, initialHeight uint64, chainID string) (types.Hash, uint64, error) {
	resp, err := c.client.InitChain(context.Background(), &pb.InitChainRequest{
		GenesisTime:   genesisTime.Unix(),
		InitialHeight: initialHeight,
		ChainId:       chainID,
	})
	if err != nil {
		return types.Hash{}, 0, err
	}

	var stateRoot types.Hash
	copy(stateRoot[:], resp.StateRoot)

	return stateRoot, resp.MaxBytes, nil
}

func (c *Client) GetTxs() ([]types.Tx, error) {
	resp, err := c.client.GetTxs(context.Background(), &pb.GetTxsRequest{})
	if err != nil {
		return nil, err
	}

	txs := make([]types.Tx, len(resp.Txs))
	for i, tx := range resp.Txs {
		txs[i] = tx
	}

	return txs, nil
}

func (c *Client) ExecuteTxs(txs []types.Tx, blockHeight uint64, timestamp time.Time, prevStateRoot types.Hash) (types.Hash, uint64, error) {
	req := &pb.ExecuteTxsRequest{
		Txs:           make([][]byte, len(txs)),
		BlockHeight:   blockHeight,
		Timestamp:     timestamp.Unix(),
		PrevStateRoot: prevStateRoot[:],
	}
	for i, tx := range txs {
		req.Txs[i] = tx
	}

	resp, err := c.client.ExecuteTxs(context.Background(), req)
	if err != nil {
		return types.Hash{}, 0, err
	}

	var updatedStateRoot types.Hash
	copy(updatedStateRoot[:], resp.UpdatedStateRoot)

	return updatedStateRoot, resp.MaxBytes, nil
}

func (c *Client) SetFinal(blockHeight uint64) error {
	_, err := c.client.SetFinal(context.Background(), &pb.SetFinalRequest{
		BlockHeight: blockHeight,
	})
	return err
}

package grpc

import (
	"context"
	"time"

	"github.com/LastL2/execution"
	pb "github.com/LastL2/types"
	"github.com/rollkit/rollkit/types"
)

// Server implements the gRPC server for execution service
type Server struct {
	pb.UnimplementedExecutionServiceServer
	exec execution.Execute
}

// NewServer creates a new gRPC server instance
func NewServer(exec execution.Execute) pb.ExecutionServiceServer {
	return &Server{
		exec: exec,
	}
}

func (s *Server) InitChain(ctx context.Context, req *pb.InitChainRequest) (*pb.InitChainResponse, error) {
	stateRoot, maxBytes, err := s.exec.InitChain(
		time.Unix(req.GenesisTime, 0),
		req.InitialHeight,
		req.ChainId,
	)
	if err != nil {
		return nil, err
	}

	return &pb.InitChainResponse{
		StateRoot: stateRoot[:],
		MaxBytes:  maxBytes,
	}, nil
}

func (s *Server) GetTxs(ctx context.Context, req *pb.GetTxsRequest) (*pb.GetTxsResponse, error) {
	txs, err := s.exec.GetTxs()
	if err != nil {
		return nil, err
	}

	pbTxs := make([][]byte, len(txs))
	for i, tx := range txs {
		pbTxs[i] = tx
	}

	return &pb.GetTxsResponse{
		Txs: pbTxs,
	}, nil
}

func (s *Server) ExecuteTxs(ctx context.Context, req *pb.ExecuteTxsRequest) (*pb.ExecuteTxsResponse, error) {
	txs := make([]types.Tx, len(req.Txs))
	for i, tx := range req.Txs {
		txs[i] = tx
	}

	var prevStateRoot types.Hash
	copy(prevStateRoot[:], req.PrevStateRoot)

	updatedStateRoot, maxBytes, err := s.exec.ExecuteTxs(
		txs,
		req.BlockHeight,
		time.Unix(req.Timestamp, 0),
		prevStateRoot,
	)
	if err != nil {
		return nil, err
	}

	return &pb.ExecuteTxsResponse{
		UpdatedStateRoot: updatedStateRoot[:],
		MaxBytes:         maxBytes,
	}, nil
}

func (s *Server) SetFinal(ctx context.Context, req *pb.SetFinalRequest) (*pb.SetFinalResponse, error) {
	err := s.exec.SetFinal(req.BlockHeight)
	if err != nil {
		return nil, err
	}

	return &pb.SetFinalResponse{}, nil
}

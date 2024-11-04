package grpc

import (
	"context"
	"time"

	"github.com/rollkit/go-execution"
	"github.com/rollkit/go-execution/types"
	pb "github.com/rollkit/go-execution/types/pb/execution"
)

type Server struct {
	pb.UnimplementedExecutionServiceServer
	exec   execution.Execute
	config *Config
}

func NewServer(exec execution.Execute, config *Config) pb.ExecutionServiceServer {
	if config == nil {
		config = DefaultConfig()
	}
	return &Server{
		exec:   exec,
		config: config,
	}
}

func (s *Server) validateAuth(ctx context.Context) error {
	if s.config.JWTSecret != nil {
		return s.validateJWT(ctx)
	}
	return nil
}

// TO-DO
func (s *Server) validateJWT(_ context.Context) error {
	return nil
}

func (s *Server) InitChain(ctx context.Context, req *pb.InitChainRequest) (*pb.InitChainResponse, error) {
	if err := s.validateAuth(ctx); err != nil {
		return nil, err
	}

	// Convert Unix timestamp to UTC time
	genesisTime := time.Unix(req.GenesisTime, 0).UTC()

	stateRoot, maxBytes, err := s.exec.InitChain(
		genesisTime,
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

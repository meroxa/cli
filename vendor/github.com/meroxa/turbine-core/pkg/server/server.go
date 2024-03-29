//go:generate mockgen -source=server.go -package=mock -destination=mock/server_mock.go Server

package server

import (
	"context"
	"log"
	"net"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	ListenAddress = "localhost:50500"
)

type Server interface {
	Run(context.Context)
	RunAddr(context.Context, string)
	GracefulStop()
}

var _ Server = (*turbineCoreServer)(nil)

type turbineCoreServer struct {
	*grpc.Server
}

func NewRunServer() *turbineCoreServer {
	s := grpc.NewServer()
	pb.RegisterTurbineServiceServer(s, NewRunService())
	return &turbineCoreServer{Server: s}
}

func NewSpecBuilderServer() *turbineCoreServer {
	s := grpc.NewServer()
	pb.RegisterTurbineServiceServer(s, NewSpecBuilderService())
	return &turbineCoreServer{Server: s}
}

func NewRecordServer() *turbineCoreServer {
	return NewSpecBuilderServer()
}

func (s *turbineCoreServer) Run(ctx context.Context) {
	s.RunAddr(ctx, ListenAddress)
}

func (s *turbineCoreServer) RunAddr(ctx context.Context, addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func empty() *emptypb.Empty {
	return new(emptypb.Empty)
}

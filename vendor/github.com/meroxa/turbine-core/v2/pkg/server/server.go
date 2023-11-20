//go:generate mockgen -source=server.go -package=mock -destination=mock/server_mock.go Server

package server

import (
	"context"
	"log"
	"net"

	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
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

var _ Server = (*TurbineCoreServer)(nil)

type TurbineCoreServer struct {
	*grpc.Server
}

func NewRunServer() *TurbineCoreServer {
	s := grpc.NewServer()
	turbinev2.RegisterServiceServer(s, NewRunService())
	return &TurbineCoreServer{Server: s}
}

func NewSpecBuilderServer() *TurbineCoreServer {
	s := grpc.NewServer()
	turbinev2.RegisterServiceServer(s, NewSpecBuilderService())
	return &TurbineCoreServer{Server: s}
}

func NewRecordServer() *TurbineCoreServer {
	return NewSpecBuilderServer()
}

func (s *TurbineCoreServer) Run(ctx context.Context) {
	s.RunAddr(ctx, ListenAddress)
}

func (s *TurbineCoreServer) RunAddr(ctx context.Context, addr string) {
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

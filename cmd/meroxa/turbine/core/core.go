package core

import (
	"fmt"
	"log"
	"net"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/servers/info"
	"github.com/meroxa/turbine-core/servers/local"
	"github.com/meroxa/turbine-core/servers/platform"
	"google.golang.org/grpc"
)

type TurbineCoreServer struct {
	server *grpc.Server
}

func NewTurbineCoreServer() TurbineCoreServer {
	return TurbineCoreServer{}
}

func (cs *TurbineCoreServer) Run(port int, mode string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	cs.server = s
	log.Printf("running gRPC server in %s mode", mode)
	switch mode {
	case "info":
		pb.RegisterTurbineServiceServer(s, info.New())
	case "platform":
		pb.RegisterTurbineServiceServer(s, platform.New())
	case "local":
		pb.RegisterTurbineServiceServer(s, local.New())
	default:
		log.Fatalf("unsupported or invalid mode %s", mode)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (cs *TurbineCoreServer) Stop() {
	cs.server.GracefulStop()
}

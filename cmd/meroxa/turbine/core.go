package turbine

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/meroxa/cli/cmd/meroxa/turbine/internal"

	"github.com/meroxa/turbine-core/pkg/client"
	"github.com/meroxa/turbine-core/pkg/server"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
)

type Core struct {
	grpcListenAddress string
	builder           server.Server
	runner            server.Server
	client            client.Client
}

func NewCore() *Core {
	return &Core{
		runner:            server.NewRunServer(),
		grpcListenAddress: internal.RandomLocalAddr(),
		builder:           server.NewSpecBuilderServer(),
	}
}

func (t *Core) Start(ctx context.Context) (string, error) {
	var (
		err     error
		retries = time.Tick(3 * time.Second)
	)

	go t.builder.RunAddr(ctx, t.grpcListenAddress)

	// NB: Spin until server is ready.
	//     Ideally, server should communicate readyness but until then we wait.
	for range retries {
		t.client, err = client.DialTimeout(t.grpcListenAddress, time.Second)
		if err == nil {
			return t.grpcListenAddress, nil
		}
	}

	return "", err
}

func (t *Core) Stop() (func(), error) {
	return func() {
		t.client.Close()
		t.builder.GracefulStop()
	}, nil
}

func (t *Core) Run(ctx context.Context) (func(), string) {
	go t.runner.RunAddr(ctx, t.grpcListenAddress)
	return t.runner.GracefulStop, t.grpcListenAddress
}

func (t *Core) GetDeploymentSpec(ctx context.Context, imageName string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.client.GetSpec(ctx, &pb.GetSpecRequest{
		Image: imageName,
	})
	if err != nil {
		return "", err
	}

	return string(resp.Spec), nil
}

// NeedsToBuild reads from the Turbine application to determine whether it needs to be built or not
// this is currently based on the number of functions.
func (t *Core) NeedsToBuild(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.client.HasFunctions(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	return resp.Value, nil
}

//go:generate mockgen -source=client.go -package=mock -destination=mock/client_mock.go Client

package client

import (
	"context"
	"time"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	Close()
	pb.TurbineServiceClient
}

type client struct {
	*grpc.ClientConn
	pb.TurbineServiceClient
}

func DialTimeout(addr string, timeout time.Duration) (*client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return DialContext(ctx, addr)
}

func DialContext(ctx context.Context, addr string) (*client, error) {
	c, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &client{
		ClientConn:           c,
		TurbineServiceClient: pb.NewTurbineServiceClient(c),
	}, nil
}

func (c *client) Close() {
	c.ClientConn.Close()
}

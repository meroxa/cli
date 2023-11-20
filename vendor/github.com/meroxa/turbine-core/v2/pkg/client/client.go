//go:generate mockgen -source=client.go -package=mock -destination=mock/client_mock.go Client

package client

import (
	"context"
	"time"

	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ Client = (*TurbineClient)(nil)

type Client interface {
	Close()
	turbinev2.ServiceClient
}

type TurbineClient struct {
	*grpc.ClientConn
	turbinev2.ServiceClient
}

func DialTimeout(addr string, timeout time.Duration) (*TurbineClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return DialContext(ctx, addr)
}

func DialContext(ctx context.Context, addr string) (*TurbineClient, error) {
	c, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &TurbineClient{
		ClientConn:    c,
		ServiceClient: turbinev2.NewServiceClient(c),
	}, nil
}

func (c *TurbineClient) Close() {
	c.ClientConn.Close()
}

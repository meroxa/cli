package turbine

import (
	"context"
)

type CLI interface {
	Upgrade(vendor bool) error
	Run(ctx context.Context) error
	Init(ctx context.Context, name string) error
	GitInit(ctx context.Context, appPath string) error
}

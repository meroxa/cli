package turbine

import (
	"context"
)

type CLI interface {
	Upgrade(vendor bool) error
	Run(ctx context.Context) error
}

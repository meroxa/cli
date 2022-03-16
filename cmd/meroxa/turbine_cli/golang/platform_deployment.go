package turbineGo

import (
	"context"

	"github.com/meroxa/cli/log"
)

// TODO: Implement
func (gd *Deploy) getPlatformImage(ctx context.Context, l log.Logger) (string, error) {
	// Get source (POST /source)
	l.Infof(ctx, "POST /source")
	// POST /builds
	l.Infof(ctx, "POST /builds")
	return "", nil
}

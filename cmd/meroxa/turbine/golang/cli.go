package turbinego

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/pkg/client"
	"github.com/meroxa/turbine-core/pkg/server"
)

type turbineGoCLI struct {
	logger            log.Logger
	appPath           string
	grpcListenAddress string

	builder server.Server
	runner  server.Server
	bc      client.Client
}

func New(l log.Logger, appPath string) turbine.CLI {
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)
	port := rand.Intn(65535-1024) + 1024 //nolint:gomnd
	serverAddress := fmt.Sprintf("localhost:%d", port)

	return &turbineGoCLI{
		logger:            l,
		appPath:           appPath,
		runner:            server.NewRunServer(),
		builder:           server.NewSpecBuilderServer(),
		grpcListenAddress: serverAddress,
	}
}

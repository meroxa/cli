package turbinejs

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/internal"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/pkg/client"
	"github.com/meroxa/turbine-core/pkg/server"
)

type turbineJsCLI struct {
	logger            log.Logger
	appPath           string
	grpcListenAddress string

	builder server.Server
	runner  server.Server
	bc      client.Client
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbineJsCLI{
		logger:            l,
		appPath:           appPath,
		runner:            server.NewRunServer(),
		builder:           server.NewSpecBuilderServer(),
		grpcListenAddress: internal.RandomLocalAddr(),
	}
}

func RunTurbineJS(ctx context.Context, l log.Logger, params ...string) (string, error) {
	args := getTurbineJSBinary(params)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	return turbine.RunCmdWithErrorDetection(ctx, cmd, l)
}

func RunTurbineJSWithEnvVars(ctx context.Context, l log.Logger, envs map[string]string, params ...string) (string, error) {
	args := getTurbineJSBinary(params)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Env = os.Environ()
	for key, value := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	return turbine.RunCmdWithErrorDetection(ctx, cmd, l)
}

func getTurbineJSBinary(params []string) []string {
	shouldUseLocalTurbineJS := global.GetLocalTurbineJSSetting()
	turbineJSBinary := fmt.Sprintf("@meroxa/turbine-js-cli@%s", TurbineJSVersion)
	if shouldUseLocalTurbineJS == "true" {
		turbineJSBinary = "turbine-js-cli"
	}
	args := []string{"npx", "--yes", turbineJSBinary}
	args = append(args, params...)
	return args
}

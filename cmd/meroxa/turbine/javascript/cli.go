package turbinejs

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

type turbineJsCLI struct {
	logger  log.Logger
	appPath string
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbineJsCLI{logger: l, appPath: appPath}
}

func RunTurbineJS(ctx context.Context, params ...string) (cmd *exec.Cmd) {
	args := getTurbineJSBinary(params)
	return executeTurbineJSCommand(ctx, args)
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

func executeTurbineJSCommand(ctx context.Context, params []string) *exec.Cmd {
	return exec.CommandContext(ctx, params[0], params[1:]...) //nolint:gosec
}

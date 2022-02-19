package turbinecli

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
)

// BuildJSApp calls turbine_cli-js to build an application.
func BuildJSApp(ctx context.Context, l log.Logger) error {
	// TODO: Handle the requirement of https://github.com/meroxa/turbine-js.git being installed
	// cd into the path first
	cmd := exec.Command("npx", "turbine_cli", "test")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	l.Info(ctx, string(stdout))
	return nil
}

func DeployJSApp(ctx context.Context, path string, l log.Logger) error {
	cmd := exec.Command("npx", "turbine_cli", "deploy", path)

	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken))

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		l.Error(ctx, string(stdout))
		return err
	}
	l.Info(ctx, string(stdout))
	return nil
}

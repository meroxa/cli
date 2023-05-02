package internal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
)

type TurbineCommand string

const (
	TurbineCommandRun     TurbineCommand = "turbine-js-run"
	TurbineCommandRecord  TurbineCommand = "turbine-js-record"
	TurbineCommandBuild   TurbineCommand = "turbine-js-dockerfile"
	TurbineCommandVersion TurbineCommand = "turbine-js-version"
)

func NewTurbineCmd(ctx context.Context, appPath string, command TurbineCommand, env map[string]string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "node", append([]string{
		path.Join(appPath, "node_modules", "@meroxa", "turbine-js", "bin", string(command)),
	}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = appPath
	cmd.Env = os.Environ()

	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	return cmd
}

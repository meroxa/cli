package turbinejs

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/meroxa/cli/cmd/meroxa/global"
	turbinecli "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
)

func NeedsToBuild(path string) (bool, error) {
	// TODO: change to `hasfunctions` when it's ready
	cmd := exec.Command("npx", "turbine", "functions", path)
	// => "true" | "false"
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(string(output))
}

// npx turbine whatever path => CLI could carry on with creating the tar.zip, post source, build...
// once that's complete, it's when we'd call `npx turbine deploy path`.
func Deploy(ctx context.Context, path string, l log.Logger) (string, error) {
	cmd := exec.Command("npx", "turbine", "deploy", path)

	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return "", err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken))

	return turbinecli.RunCmdWithErrorDetection(ctx, cmd, l)
}

// TODO: Have a script to cleanup the temp directory used (right after source is uploaded)

// TODO: Add a function that creates the needed structure for a JS app (separate from the deploy step)

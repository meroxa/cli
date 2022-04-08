package turbinejs

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/global"
	turbinecli "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
)

func NeedsToBuild(path string) (bool, error) {
	cmd := exec.Command("npx", "turbine", "hasfunctions", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		err := fmt.Errorf(
			"unable to determine if the Meroxa Application at %s has a Process; %s",
			path,
			string(output))
		return false, err
	}
	return strconv.ParseBool(strings.TrimSpace(string(output)))
}

func BuildApp(path string) (string, error) {
	cmd := exec.Command("npx", "turbine", "clibuild", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("unable to build Meroxa Application at %s; %s", path, string(output))
	}
	return strings.TrimSpace(string(output)), err
}

func Deploy(ctx context.Context, path string, l log.Logger) (string, error) {
	cmd := exec.Command("npx", "turbine", "clideploy", path)

	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return "", err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken))

	return turbinecli.RunCmdWithErrorDetection(ctx, cmd, l)
}

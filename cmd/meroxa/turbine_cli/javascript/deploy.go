package turbinejs

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/meroxa/cli/cmd/meroxa/global"
	turbinecli "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
)

func NeedsToBuild(path string) (bool, error) {
	cmd := exec.Command("npx", "--yes", "@meroxa/turbine-js@0.2.0", "hasfunctions", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		err := fmt.Errorf(
			"unable to determine if the Meroxa Application at %s has a Process; %s",
			path,
			string(output))
		return false, err
	}

	r := regexp.MustCompile("\nturbine-response: (true|false)\n")
	match := r.FindStringSubmatch(string(output))
	if match == nil || len(match) < 2 {
		err := fmt.Errorf(
			"unable to determine if the Meroxa Application at %s has a Process; %s",
			path,
			string(output))
		return false, err
	}
	return strconv.ParseBool(match[1])
}

func RunDeployApp(ctx context.Context, l log.Logger, path, imageName string) (string, error) {
	cmd := exec.Command("npx", "--yes", "@meroxa/turbine-js@0.2.0", "deploy", imageName, path)

	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return "", err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken))

	return turbinecli.RunCmdWithErrorDetection(ctx, cmd, l)
}

func CreateDockerfile(path string) error {
	cmd := exec.Command("npx", "--yes", "@meroxa/turbine-js@0.2.0", "generatedockerfile", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to build Meroxa Application at %s; %s", path, string(output))
	}

	return nil
}

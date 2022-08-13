package turbinejs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

func NeedsToBuild(path string) (bool, error) {
	cmd := RunTurbineJS("hasfunctions", path)
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

func CreateDockerfile(ctx context.Context, l log.Logger, path string) error {
	cmd := RunTurbineJS("clibuild", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to create Dockerfile at %s; %s", path, string(output))
	}

	return err
}

func RunDeployApp(ctx context.Context, l log.Logger, path, imageName, appName, gitSha, specVersion string) error {
	params := []string{"clideploy", imageName, path, appName, gitSha}

	if specVersion != "" {
		params = append(params, specVersion)
	}

	cmd := RunTurbineJS(params...)

	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken))

	_, err = turbine.RunCmdWithErrorDetection(ctx, cmd, l)
	return err
}

// GetResources asks turbine for a list of resources used by the given app.
func GetResources(ctx context.Context, l log.Logger, appPath, appName string) ([]turbine.ApplicationResource, error) {
	var resources []turbine.ApplicationResource

	cmd := RunTurbineJS("listresources", appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return resources, errors.New(string(output))
	}

	if err := json.Unmarshal(output, &resources); err != nil {
		// fall back if not json
		return turbine.GetResourceNamesFromString(string(output)), nil
	}

	return resources, nil
}

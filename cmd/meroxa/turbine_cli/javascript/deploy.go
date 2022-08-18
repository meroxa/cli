package turbinejs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/global"
	turbinecli "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
)

func NeedsToBuild(path string) (bool, error) {
	cmd := turbinecli.RunTurbineJS("hasfunctions", path)
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
	cmd := turbinecli.RunTurbineJS("clibuild", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to create Dockerfile at %s; %s", path, string(output))
	}

	return err
}

func RunDeployApp(ctx context.Context, l log.Logger, path, imageName, appName, gitSha string) error {
	cmd := turbinecli.RunTurbineJS("clideploy", imageName, path, appName, gitSha)

	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken))

	_, err = turbinecli.RunCmdWithErrorDetection(ctx, cmd, l)
	return err
}

// GetResourceNames asks turbine for a list of resources used by the given app.
func GetResourceNames(ctx context.Context, l log.Logger, appPath, appName string) ([]string, error) {
	var names []string

	cmd := turbinecli.RunTurbineJS("listresources", appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return names, errors.New(string(output))
	}

	var parsed []turbinecli.ApplicationResource
	if err := json.Unmarshal(output, &parsed); err != nil {
		return getResourceNamesFromString(string(output)), nil
	}

	for i := range parsed {
		kv := parsed[i]
		if kv.Name != "" {
			names = append(names, kv.Name)
		}
	}

	return names, nil
}

func getResourceNamesFromString(s string) []string {
	var names []string

	r := regexp.MustCompile(`\[(.+?)\]`)
	sliceString := r.FindStringSubmatch(s)
	if len(sliceString) > 0 {
		names = strings.Fields(sliceString[1])
	}

	return names
}

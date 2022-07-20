package turbinejs

import (
	"context"
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

func BuildApp(path string) (string, error) {
	cmd := turbinecli.RunTurbineJS("clibuild", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("unable to build Meroxa Application at %s; %s", path, string(output))
	}

	r := regexp.MustCompile("\nturbine-response: (.*)\n")
	match := r.FindStringSubmatch(string(output))
	if match == nil || len(match) < 2 {
		return "", fmt.Errorf("unable to build Meroxa Application at %s; %s", path, string(output))
	}
	return match[1], err
}

func RunDeployApp(ctx context.Context, l log.Logger, path, imageName, gitSha string) error {
	cmd := turbinecli.RunTurbineJS("clideploy", imageName, path, gitSha)

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
	r := regexp.MustCompile("\nturbine-response: (.*)\n")
	match := r.FindStringSubmatch(string(output))
	if match == nil || len(match) < 2 {
		return names, fmt.Errorf("unable to verify resource availability for Meroxa Application at %s; %s", appPath, string(output))
	}
	text := match[1]
	names = strings.Split(text, ",")
	for i, name := range names {
		names[i] = strings.TrimSpace(name)
	}

	return names, nil
}

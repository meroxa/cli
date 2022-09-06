package turbinepy

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/meroxa/cli/log"
)

// Build created the needed structure for a python app.
func (t *turbinePyCLI) Build(ctx context.Context, appName string, platform bool) (string, error) {
	t.logger.StartSpinner("\t", " Building application...")
	cmd := exec.CommandContext(ctx, "turbine-py", "clibuild", t.appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.logger.StopSpinnerWithStatus("\t", log.Failed)
		return "", fmt.Errorf("unable to build Meroxa Application at %s; %s", t.appPath, string(output))
	}
	r := regexp.MustCompile("^turbine-response: ([^\n]*)")
	match := r.FindStringSubmatch(string(output))
	if match == nil || len(match) < 2 {
		t.logger.StopSpinnerWithStatus("\t", log.Failed)
		return "", fmt.Errorf("unable to build Meroxa Application at %s; %s", t.appPath, string(output))
	}
	t.logger.StopSpinnerWithStatus("Application built", log.Successful)
	return match[1], err
}

func (t *turbinePyCLI) CleanUpBinaries(ctx context.Context, appName string) {
}

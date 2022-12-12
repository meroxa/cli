package turbinego

import (
	"fmt"
	"os/exec"

	"github.com/meroxa/cli/log"
)

// GoGetDeps updates the latest Meroxa mods.
func GoGetDeps(l log.Logger) error {
	l.StartSpinner("\t", "Getting latest turbine-go dependency...")
	cmd := exec.Command("go", "get", "-u", "github.com/meroxa/turbine-go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		l.StopSpinnerWithStatus(fmt.Sprintf("%s", string(output)), log.Failed)
		return err
	}

	l.StopSpinnerWithStatus("Downloaded latest turbine-go dependency successfully!", log.Successful)
	return nil
}

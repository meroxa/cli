package turbinego

import (
	"fmt"
	"os"
	"os/exec"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

// Upgrade fetches the latest Meroxa dependencies.
func (t *turbineGoCLI) Upgrade(vendor bool) error {
	pwd, err := utils.SwitchToAppDirectory(t.appPath)
	if err != nil {
		return err
	}

	err = GoGetDeps(t.logger)
	if err != nil {
		return err
	}

	err = t.tidy(vendor)
	if err != nil {
		return err
	}

	return os.Chdir(pwd)
}

func (t *turbineGoCLI) tidy(vendor bool) error {
	var err error
	t.logger.StartSpinner("\t", " Tidying up Golang modules...")
	if vendor {
		_, err = os.Stat("vendor")
		if !os.IsNotExist(err) {
			cmd := exec.Command("go", "mod", "vendor")
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.logger.StopSpinnerWithStatus("Failed to download modules to vendor", log.Failed)
				return fmt.Errorf(string(output))
			}
		}
		cmd := exec.Command("go", "mod", "tidy")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.logger.StopSpinnerWithStatus("Failed to tidy modules", log.Failed)
			return fmt.Errorf(string(output))
		}
	}
	t.logger.StopSpinnerWithStatus("Finished tidying up Golang modules successfully!", log.Successful)
	return nil
}

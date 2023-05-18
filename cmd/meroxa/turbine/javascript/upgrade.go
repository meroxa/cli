package turbinejs

import (
	"fmt"
	"os/exec"

	"github.com/meroxa/cli/log"
)

// Upgrade fetches the latest Meroxa dependencies.
func (t *turbineJsCLI) Upgrade(_ bool) error {
	cmd := exec.Command("grep", "turbine-js", "package.json")
	cmd.Dir = t.appPath
	err := cmd.Wait()
	if err != nil {
		t.logger.StartSpinner("\t", "Adding @meroxa/turbine-js requirement...")
		cmd = exec.Command("npm", "install", "@meroxa/turbine-js", "--save")
		cmd.Dir = t.appPath
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.logger.StopSpinnerWithStatus("Failed to install @meroxa/turbine-js", log.Failed)
			return fmt.Errorf(string(out))
		}

		t.logger.StopSpinnerWithStatus("Added @meroxa/turbine-js requirement successfully!", log.Successful)
	} else {
		t.logger.StartSpinner("\t", "Upgrading @meroxa/turbine-js...")
		cmd = exec.Command("npm", "upgrade", "@meroxa/turbine-js")
		cmd.Dir = t.appPath
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.logger.StopSpinnerWithStatus("Failed to upgrade @meroxa/turbine-js", log.Failed)
			return fmt.Errorf(string(out))
		}
		t.logger.StopSpinnerWithStatus("Upgraded @meroxa/turbine-js!", log.Successful)
	}
	return nil
}

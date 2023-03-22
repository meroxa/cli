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
		t.logger.StartSpinner("\t", "Adding @meroxa/turbine-js-framework requirement...")
		cmd = exec.Command("npm", "install", "@meroxa/turbine-js-framework", "--save")
		cmd.Dir = t.appPath
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.logger.StopSpinnerWithStatus("Failed to install @meroxa/turbine-js-framework", log.Failed)
			return fmt.Errorf(string(out))
		}

		cmd = exec.Command("npm", "uninstall", "@meroxa/turbine-js", "--save")
		cmd.Dir = t.appPath
		err = cmd.Run()
		if err != nil {
			t.logger.StopSpinnerWithStatus("Failed to uninstall @meroxa/turbine-js. Moving on...", log.Failed)
		}

		cmd = exec.Command("npm", "update")
		cmd.Dir = t.appPath
		out, err = cmd.CombinedOutput()
		if err != nil {
			t.logger.StopSpinnerWithStatus("Failed to run npm update", log.Failed)
			return fmt.Errorf(string(out))
		}
		t.logger.StopSpinnerWithStatus("Added @meroxa/turbine-js-framework requirement successfully!", log.Successful)
	} else {
		t.logger.StartSpinner("\t", "Upgrading @meroxa/turbine-js-framework...")
		cmd = exec.Command("npm", "upgrade", "@meroxa/turbine-js-framework")
		cmd.Dir = t.appPath
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.logger.StopSpinnerWithStatus("Failed to upgrade @meroxa/turbine-js-framework", log.Failed)
			return fmt.Errorf(string(out))
		}
		t.logger.StopSpinnerWithStatus("Upgraded @meroxa/turbine-js-framework successfully!", log.Successful)
	}
	return nil
}

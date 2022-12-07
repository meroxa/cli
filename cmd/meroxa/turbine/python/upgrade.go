package turbinepy

import (
	"fmt"
	"os/exec"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

const turbinePYVersion = "1.5.3"

// Upgrade fetches the latest Meroxa dependencies.
func (t *turbinePyCLI) Upgrade(vendor bool) error {
	t.logger.StartSpinner("\t", "Upgrading turbine dependencies ")

	// Run pip upgrade
	cmd := exec.Command("pip", "install", "-U", "turbine-py", "meroxa-py")
	cmd.Dir = t.appPath
	err1 := cmd.Run()

	// Get current version from newly updated turbine
	cmd = exec.Command("turbine-py", "version")
	output, err2 := cmd.Output()
	turbineVersion, _ := utils.GetTurbineResponseFromOutput(string(output))

	replace := fmt.Sprintf("'s+turbine-py+turbine-py==%v+g'", turbineVersion)
	fmt.Printf("The replace is %s\n", replace)

	cmd = exec.Command("bash", "-c", "sed -i "+replace+" requirements.txt")
	cmd.Dir = t.appPath
	err4 := cmd.Run()

	t.logger.StopSpinnerWithStatus("Tidied up requirements.txt successfully!", log.Successful)

	if err1 == nil && err2 == nil && err4 == nil {
		t.logger.StopSpinnerWithStatus("Tidied up requirements.txt successfully!", log.Successful)
	} else {
		t.logger.StopSpinnerWithStatus("Issues encountered tidying up requirements.txt. Moving on...", log.Failed)
	}

	t.logger.StopSpinnerWithStatus("Updated Python dependencies successfully!", log.Successful)
	return nil
}

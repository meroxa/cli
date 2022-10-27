package turbinepy

import (
	"fmt"
	"os/exec"

	"github.com/meroxa/cli/log"
)

const turbinePYVersion = "1.5.3"

// Upgrade fetches the latest Meroxa dependencies.
func (t *turbinePyCLI) Upgrade(vendor bool) error {
	cmd := exec.Command("grep", "turbine-py==", "requirements.txt")
	cmd.Dir = t.appPath
	err := cmd.Run()
	if err != nil {
		t.logger.StartSpinner("\t", "Tidying up requirements.txt...")
		cmd = exec.Command("bash", "-c", "sed -i 's+meroxa-py++g' requirements.txt")
		cmd.Dir = t.appPath
		err1 := cmd.Run()

		replace := fmt.Sprintf("'s+turbine-py+turbine-py==%s+g'", turbinePYVersion)
		cmd = exec.Command("bash", "-c", "sed -i "+replace+" requirements.txt")
		cmd.Dir = t.appPath
		err2 := cmd.Run()
		if err1 == nil && err2 == nil {
			t.logger.StopSpinnerWithStatus("Tidied up requirements.txt successfully!", log.Successful)
		} else {
			t.logger.StopSpinnerWithStatus("Issues encountered tidying up requirements.txt. Moving on...", log.Failed)
		}
	}

	t.logger.StartSpinner("\t", "Updating Python dependencies...")
	cmd = exec.Command("pip", "install", "turbine-py", "-U")
	cmd.Dir = t.appPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.logger.StopSpinnerWithStatus("\t", log.Failed)
		return fmt.Errorf(string(out))
	}
	t.logger.StopSpinnerWithStatus("Updated Python dependencies successfully!", log.Successful)
	return nil
}

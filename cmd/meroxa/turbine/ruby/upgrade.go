package turbinerb

import (
	"os/exec"
)

// Upgrade fetches the latest Meroxa dependencies.
func (t *turbineRbCLI) Upgrade(vendor bool) error {
	cmd := exec.Command("ruby command goes here")
	cmd.Dir = t.appPath
	err := cmd.Wait()
	if err != nil {
		return err
	}
	//upgrade stub for ruby
	return nil
}

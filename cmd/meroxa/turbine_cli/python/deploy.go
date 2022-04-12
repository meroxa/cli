package turbinepy

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// TODO: Add a function that creates the needed structure for a python app

// NeedsToBuild determines if the app has functions or not.
func NeedsToBuild(path string) (bool, error) {
	cmd := exec.Command("turbine-py", "hasFunctions", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		err := fmt.Errorf(
			"unable to determine if the Meroxa Application at %s has a Process; %s",
			path,
			string(output))
		return false, err
	}
	return strconv.ParseBool(strings.TrimSpace(string(output)))
}

// TODO: Add a function that actually creates the meroxa resources...

// TODO: Have a script to cleanup the temp directory used (right after source is uploaded)

package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

type TurbineCommand string

var (
	TurbineCommandRun   TurbineCommand = "TurbineRb.run"
	TurbineCommandBuild TurbineCommand = "TurbineRb.build"
)

func NewTurbineCmd(appPath string, command TurbineCommand, env map[string]string) *exec.Cmd {
	cmd := exec.Command("ruby", []string{
		"-r", path.Join(appPath, "app"),
		"-e", string(command),
	}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = appPath
	cmd.Env = os.Environ()

	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	return cmd
}

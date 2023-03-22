package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

type TurbineCommand string

const (
	TurbineCommandRun    TurbineCommand = "TurbineRb.run"
	TurbineCommandRecord TurbineCommand = "TurbineRb.record"
	TurbineCommandBuild  TurbineCommand = "TurbineRb.build"
)

func NewTurbineCmd(appPath string, command TurbineCommand, env map[string]string, flags ...string) *exec.Cmd {
	cmd := exec.Command("ruby", append([]string{
		"-r", path.Join(appPath, "app"),
		"-e", string(command)}, flags...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = appPath
	cmd.Env = os.Environ()

	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	return cmd
}

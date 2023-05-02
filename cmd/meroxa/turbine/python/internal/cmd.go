package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TurbineCommand string

const (
	TurbineCommandRun     TurbineCommand = "run"
	TurbineCommandRecord  TurbineCommand = "record"
	TurbineCommandBuild   TurbineCommand = "build"
	TurbineCommandVersion TurbineCommand = `version`
)

func RunTurbineCmd(appPath string, command TurbineCommand, env map[string]string, flags ...string) (string, error) {
	cmd := NewTurbineCmd(appPath, command, env, flags...)
	cmd.Stderr = nil
	cmd.Stdout = nil
	cmd.Dir = appPath
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func NewTurbineCmd(appPath string, command TurbineCommand, env map[string]string, flags ...string) *exec.Cmd {
	cmd := exec.Command("turbine-py", append([]string{
		string(command),
	}, flags...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = appPath
	cmd.Env = os.Environ()

	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	return cmd
}

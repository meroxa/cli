package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type TurbineCommand string

const (
	TurbineCommandRun    TurbineCommand = "TurbineRb.run"
	TurbineCommandRecord TurbineCommand = "TurbineRb.record"
	TurbineCommandBuild  TurbineCommand = "TurbineRb.build"
	TurbineCommandVersion TurbineCommand = `puts Gem.loaded_specs["turbine_rb"].version.version`
)

func RunTurbineCmd(appPath string, command TurbineCommand, env map[string]string, flags ...string) (string, error) {
	cmd := NewTurbineCmd(appPath, command, env, flags...)
	cmd.Stderr = nil
	cmd.Stdout = nil
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func NewTurbineCmd(appPath string, command TurbineCommand, env map[string]string, flags ...string) *exec.Cmd {
	cmd := exec.Command("ruby", append([]string{
		"-I", appPath,
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

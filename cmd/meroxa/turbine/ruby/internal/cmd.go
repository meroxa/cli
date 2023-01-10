package internal

import (
	"os/exec"
	"path"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
)

type TurbineRbCommand string

var (
	TurbineCommandRun    TurbineRbCommand = "TurbineRb.run"
	TurbineCommandRecord TurbineRbCommand = "TurbineRb.record"
	TurbineCommandBuild  TurbineRbCommand = "TurbineRb.build"
)

func NewTurbineRbCmd(appPath string, command TurbineRbCommand, env map[string]string, flags ...string) *exec.Cmd {
	cmd := exec.Command("ruby", append([]string{
		"-r", path.Join(appPath, "app"),
		"-e", string(command)}, flags...)...)
	return turbine.NewTurbineCmd(cmd, appPath, env)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// cmd.Dir = appPath
	// cmd.Env = os.Environ()

	// for k, v := range env {
	// 	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	// }
	// return cmd
}

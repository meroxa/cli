package turbinejs

import (
	"fmt"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/global"
)

const turbineJSVersion = "1.0.0"

var isTrue = "true"

func RunTurbineJS(params ...string) (cmd *exec.Cmd) {
	args := getTurbineJSBinary(params)
	return executeTurbineJSCommand(args)
}

func getTurbineJSBinary(params []string) []string {
	shouldUseLocalTurbineJS := global.GetLocalTurbineJSSetting()
	turbineJSBinary := fmt.Sprintf("@meroxa/turbine-js-cli@%s", turbineJSVersion)
	if shouldUseLocalTurbineJS == isTrue {
		turbineJSBinary = "turbine-js-cli"
	}
	args := []string{"npx", "--yes", turbineJSBinary}
	args = append(args, params...)
	return args
}

func executeTurbineJSCommand(params []string) *exec.Cmd {
	return exec.Command(params[0], params[1:]...) //nolint:gosec
}

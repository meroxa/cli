package turbinejs

import (
	"context"
	"os"
	"os/exec"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/turbine-core/pkg/app"
	"github.com/meroxa/turbine-core/pkg/ir"
)

func (t *turbineJsCLI) Init(_ context.Context, appName string) error {
	err := app.NewAppInit(appName, ir.JavaScript, t.appPath).Init()
	if err != nil {
		return err
	}

	err = jsInit(t.appPath + "/" + appName)

	if err != nil {
		return err
	}

	return nil
}

func jsInit(appPath string) error {
	// temporarily switching to the app's directory
	pwd, err := utils.SwitchToAppDirectory(appPath)
	if err != nil {
		return err
	}

	cmd := exec.Command("npm", "install")
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return os.Chdir(pwd)
}

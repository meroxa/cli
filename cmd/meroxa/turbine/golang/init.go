package turbinego

import (
	"context"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/meroxa/turbine-core/pkg/app"
	"github.com/meroxa/turbine-core/pkg/ir"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

func (t *turbineGoCLI) Init(_ context.Context, appName string) error {
	return app.NewAppInit(appName, ir.GoLang, t.appPath).Init()
}

func (t *turbineGoCLI) GitInit(ctx context.Context, name string) error {
	return utils.GitInit(ctx, name)
}

func GoInit(l log.Logger, appPath string, skipInit, vendor bool) error {
	l.StartSpinner("\t", "Running golang module initializing...")
	skipLog := "skipping go module initialization\n\tFor guidance, visit " +
		"https://docs.meroxa.com/beta-overview#go-mod-init-for-a-new-golang-turbine-data-application"
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = build.Default.GOPATH
	}
	if goPath == "" {
		l.StopSpinnerWithStatus("$GOPATH not set up; "+skipLog, log.Warning)
		return nil
	}
	i := strings.Index(appPath, filepath.Join(goPath, "src"))
	if i == -1 || i != 0 {
		l.StopSpinnerWithStatus(fmt.Sprintf("%s is not under $GOPATH/src; %s", appPath, skipLog), log.Warning)
		return nil
	}

	// temporarily switching to the app's directory
	pwd, err := utils.SwitchToAppDirectory(appPath)
	if err != nil {
		l.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	// initialize the user's module
	err = utils.SetModuleInitInAppJSON(appPath, skipInit)
	if err != nil {
		l.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	err = modulesInit(l, appPath, skipInit, vendor)
	if err != nil {
		l.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	return os.Chdir(pwd)
}

func modulesInit(l log.Logger, appPath string, skipInit, vendor bool) error {
	if skipInit {
		return nil
	}

	cmd := exec.Command("go", "mod", "init")
	output, err := cmd.CombinedOutput()
	if err != nil {
		l.StopSpinnerWithStatus(fmt.Sprintf("%s", string(output)), log.Failed)
		return err
	}
	successLog := "go mod init succeeded"
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		successLog += fmt.Sprintf(" (while assuming GOPATH is %s)", build.Default.GOPATH)
	}
	l.StopSpinnerWithStatus(successLog+"!", log.Successful)

	err = GoGetDeps(l)
	if err != nil {
		return err
	}

	// download dependencies
	err = utils.SetVendorInAppJSON(appPath, vendor)
	if err != nil {
		return err
	}

	// tidy
	goTidy := exec.Command("go", "mod", "tidy")
	if _, err = goTidy.CombinedOutput(); err != nil {
		return err
	}

	depsLog := "Downloading dependencies"
	cmd = exec.Command("go", "mod", "download")
	if vendor {
		depsLog += " to vendor"
		cmd = exec.Command("go", "mod", "vendor")
	}
	depsLog += "..."
	l.StartSpinner("\t", depsLog)
	output, err = cmd.CombinedOutput()
	if err != nil {
		l.StopSpinnerWithStatus(fmt.Sprintf("%s", string(output)), log.Failed)
		return err
	}
	l.StopSpinnerWithStatus("Downloaded all other dependencies successfully!", log.Successful)
	return nil
}

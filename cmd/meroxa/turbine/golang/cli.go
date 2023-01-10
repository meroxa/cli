package turbinego

import (
	"context"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/pkg/server"
)

type TurbineGoCommand string

var (
	TurbineCommandRun    TurbineGoCommand = "--run"
	TurbineCommandRecord TurbineGoCommand = "--record"
	TurbineCommandBuild  TurbineGoCommand = "--build"
)

type turbineGoCLI struct {
	logger            log.Logger
	appName           string
	appPath           string
	runServer         turbine.TurbineServer
	recordServer      turbine.TurbineServer
	recordClient      turbine.RecordClient
	grpcListenAddress string
}

func New(l log.Logger, appPath string) (turbine.CLI, error) {
	appName, err := utils.GetAppNameFromAppJSON(context.TODO(), l, appPath)
	if err != nil {
		return nil, err
	}

	return &turbineGoCLI{
		appName:           appName,
		logger:            l,
		appPath:           appPath,
		runServer:         server.NewRunServer(),
		recordServer:      server.NewRecordServer(),
		grpcListenAddress: server.ListenAddress,
	}, nil
}

func NewTurbineGoCmd(appPath, appName string, command TurbineGoCommand, env map[string]string, args ...string) *exec.Cmd {
	cmd := exec.Command(appName, append([]string{string(command)}, args...)...)
	return turbine.NewTurbineCmd(cmd, appPath, env)
}

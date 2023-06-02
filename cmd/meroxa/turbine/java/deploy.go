package turbinejava

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
)

//go:embed templates
var templates embed.FS

func (t *turbineJavaCLI) CreateDockerfile(_ context.Context, _ string) (string, error) {
	f, err := os.Create(path.Join(t.appPath, "Dockerfile"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	tpl, err := template.ParseFS(templates, "templates/Dockerfile.tpl")
	if err != nil {
		return "", err
	}
	if err := tpl.Execute(f, map[string]string{}); err != nil {
		return "", err
	}

	return t.appPath, nil
}

func (t *turbineJavaCLI) StartGrpcServer(ctx context.Context, _ string) (func(), error) {
	grpcListenAddress, err := t.Core.Start(ctx)
	if err != nil {
		return nil, err
	}

	if err = t.build(ctx); err != nil {
		return nil, fmt.Errorf("couldn't build Java app: %w", err)
	}

	var cmd *exec.Cmd

	cmd, err = t.run(ctx, grpcListenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed building run command: %w", err)
	}

	if _, err := turbine.RunCmdWithErrorDetection(ctx, cmd, t.logger); err != nil {
		return nil, err
	}

	return t.Core.Stop()
}

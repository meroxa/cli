package turbinego

import (
	"context"
	"embed"
	"os"
	"os/exec"
	"path"
	"text/template"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
)

//go:embed templates
var templates embed.FS

func (t *turbineGoCLI) CreateDockerfile(_ context.Context, appName string) (string, error) {
	f, err := os.Create(path.Join(t.appPath, "Dockerfile"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	tpl, err := template.ParseFS(templates, "templates/Dockerfile.tpl")
	if err != nil {
		return "", err
	}
	if err := tpl.Execute(f, map[string]string{
		"GoVersion": "1.20",
		"AppName":   appName,
	}); err != nil {
		return "", err
	}

	return t.appPath, nil
}

func (t *turbineGoCLI) StartGrpcServer(ctx context.Context, gitSha string) (func(), error) {
	grpcListenAddress, err := t.Core.Start(ctx)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("go", []string{
		"run",
		"./...",
		"build",
		"-gitsha",
		gitSha,
		"-turbine-core-server", grpcListenAddress,
		"-app-path", t.appPath,
	}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = t.appPath
	cmd.Env = os.Environ()

	if err := turbine.RunCMD(ctx, t.logger, cmd); err != nil {
		return nil, err
	}

	return t.Core.Stop()
}

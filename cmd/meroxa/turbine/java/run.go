package turbinejava

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *turbineJavaCLI) Run(ctx context.Context) error {
	gracefulStop, grpcListenAddress := t.Core.Run(ctx)
	defer gracefulStop()

	if err := t.build(ctx); err != nil {
		return fmt.Errorf("couldn't build Java app: %w", err)
	}

	runCMD, err := t.runCMD(ctx, grpcListenAddress)
	if err != nil {
		return fmt.Errorf("failed building run command: %w", err)
	}

	return turbine.RunCMD(ctx, t.logger, runCMD)
}

func (t *turbineJavaCLI) build(ctx context.Context) error {
	cmd := exec.Command(
		"mvn",
		"-U",
		"clean",
		"package",
		"-Dquarkus.package.type=uber-jar",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = t.appPath
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		t.logger.Errorf(ctx, err.Error())
		return err
	}

	return nil
}

func (t *turbineJavaCLI) runCMD(ctx context.Context, grpcAddress string) (*exec.Cmd, error) {
	jarPath, err := t.getJARPath()
	if err != nil {
		return nil, fmt.Errorf("couldn't get path to JAR file: %w", err)
	}
	cmd := exec.CommandContext(
		ctx,
		"java",
		"-jar",
		"-Dturbine.core.server="+grpcAddress,
		"-Dturbine.app.path="+t.appPath,
		"-Dturbine.mode=local",
		jarPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = t.appPath
	cmd.Env = os.Environ()

	return cmd, nil
}

func (t *turbineJavaCLI) getJARPath() (string, error) {
	targetDir := filepath.Join(t.appPath, "target")

	// Read the directory entries in the target directory
	files, err := os.ReadDir(targetDir)
	if err != nil {
		return "", fmt.Errorf("couldn't read directory %v: %w", targetDir, err)
	}

	// Iterate over the directory entries and find the file ending in "-runner.jar"
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "-runner.jar") {
			// Build the path to the file
			filePath := filepath.Join(targetDir, file.Name())
			return filePath, nil
		}
	}

	return "", fmt.Errorf("file not found in the directory %v", targetDir)
}

package turbinerb

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
)

func (t *turbineRbCLI) Run(ctx context.Context) error {
	go t.runServer.Run(ctx)
	defer t.runServer.GracefulStop()

	cmd := exec.Command("ruby", []string{
		"-r", path.Join(t.appPath, "app"),
		"-e", "Turbine.run",
	}...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("TURBINE_CORE_SERVER=%s", t.grpcListenAddress))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = t.appPath

	if err := cmd.Start(); err != nil {
		t.logger.Errorf(ctx, err.Error())
		return err
	}

	if err := cmd.Wait(); err != nil {
		t.logger.Errorf(ctx, err.Error())
		return err
	}

	return nil
}

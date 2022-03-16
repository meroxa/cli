package turbinejs

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
)

func Deploy(ctx context.Context, path string, l log.Logger) error {
	cmd := exec.Command("npx", "turbine", "deploy", path)

	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken))

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		l.Error(ctx, string(stdout))
		return err
	}
	l.Info(ctx, string(stdout))
	return nil
}

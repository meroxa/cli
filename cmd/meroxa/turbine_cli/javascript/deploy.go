package turbinejs

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/global"
	turbinecli "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
)

// npx turbine whatever path => CLI could carry on with creating the tar.zip, post source, build...
// once that's complete, it's when we'd call `npx turbine deploy path`.
func Deploy(ctx context.Context, path string, l log.Logger) error {
	cmd := exec.Command("npx", "turbine", "deploy", path)

	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken))

	_, err = turbinecli.RunCmdWithErrorDetection(ctx, cmd, l)
	return err
}


// 1. we build binary (Go) // we set up the app structure (JS/Python) <- CLI could do this
// 2. we create the docker file (each turbine-lib does this)
// 3. we call the binary passing --platform ("deploying") //
package internal_test

import (
	"context"
	"os"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/turbine/ruby/internal"
	"github.com/stretchr/testify/assert"
)

func Test_NewTurbineCmd(t *testing.T) {
	ctx := context.Background()
	cmd := internal.NewTurbineCmd(
		ctx,
		"/my/path",
		internal.TurbineCommandRun,
		map[string]string{
			"foo": "bar",
			"x":   "y",
		})

	assert.Equal(t, cmd.Stdout, os.Stdout)
	assert.Equal(t, cmd.Stderr, os.Stderr)
	assert.Equal(t, cmd.Dir, "/my/path")
	assert.Equal(t, cmd.Args, []string{
		"ruby",
		"-r", "/my/path/app",
		"-e", "TurbineRb.run",
	})
	assert.Contains(t, cmd.Env, "foo=bar")
	assert.Contains(t, cmd.Env, "x=y")
}

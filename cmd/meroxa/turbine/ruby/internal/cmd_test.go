package internal_test

import (
	"os"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/turbine/ruby/internal"
	"github.com/stretchr/testify/assert"
)

func Test_NewTurbineCmd(t *testing.T) {
	cmd := internal.NewTurbineCmd("/my/path", map[string]string{
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

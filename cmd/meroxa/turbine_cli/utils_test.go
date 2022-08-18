package turbinecli

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTurbineJSBinary(t *testing.T) {
	testCases := []struct {
		name    string
		envVar  string
		wantCmd string
	}{
		{
			name:    "MEROXA_USE_LOCAL_TURBINE_JS is unset",
			envVar:  "",
			wantCmd: fmt.Sprintf("@meroxa/turbine-js-cli@%s", turbineJSVersion),
		},
		{
			name:    "MEROXA_USE_LOCAL_TURBINE_JS is true",
			envVar:  "true",
			wantCmd: "turbine-js-cli",
		},
		{
			name:    "MEROXA_USE_LOCAL_TURBINE_JS is false",
			envVar:  "false",
			wantCmd: fmt.Sprintf("@meroxa/turbine-js-cli@%s", turbineJSVersion),
		},
		{
			name:    "MEROXA_USE_LOCAL_TURBINE_JS is set to a value that is neither true nor false",
			envVar:  "jam",
			wantCmd: fmt.Sprintf("@meroxa/turbine-js-cli@%s", turbineJSVersion),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("MEROXA_USE_LOCAL_TURBINE_JS", tc.envVar)

			params := []string{"foo", "bar"}
			result := getTurbineJSBinary(params)
			assert.Equal(t, []string{"npx", "--yes", tc.wantCmd, "foo", "bar"}, result)
		})
	}
}

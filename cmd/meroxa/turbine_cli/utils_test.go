package turbinecli

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			wantCmd: fmt.Sprintf("@meroxa/turbine-js@%s", turbineJSVersion),
		},
		{
			name:    "MEROXA_USE_LOCAL_TURBINE_JS is true",
			envVar:  "true",
			wantCmd: "turbine",
		},
		{
			name:    "MEROXA_USE_LOCAL_TURBINE_JS is false",
			envVar:  "false",
			wantCmd: fmt.Sprintf("@meroxa/turbine-js@%s", turbineJSVersion),
		},
		{
			name:    "MEROXA_USE_LOCAL_TURBINE_JS is set to a value that is neither true nor false",
			envVar:  "jam",
			wantCmd: fmt.Sprintf("@meroxa/turbine-js@%s", turbineJSVersion),
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

func TestGetPipelineUUID(t *testing.T) {
	uuid := uuid2.New().String()

	testCases := []struct {
		desc string
		logs string
		err  error
	}{
		{
			desc: "Find pipeline when app has underscores",
			logs: fmt.Sprintf(`hey\npipeline: "turbine-pipeline-n_a_m_e" (%s)\nhello`, uuid),
			err:  nil,
		},
		{
			desc: "Find pipeline when app has dashes",
			logs: fmt.Sprintf(`hey\npipeline: "turbine-pipeline-n-a-m-e" (%s)\nhello`, uuid),
			err:  nil,
		},
		{
			desc: "Fail to find pipeline when UUID is missing",
			logs: `hey\npipeline: "turbine-pipeline-n-a-m-e" \nhello`,
			err:  fmt.Errorf("pipeline UUID not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result, err := GetPipelineUUID(tc.logs)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, err, tc.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, result, uuid)
			}
		})
	}
}

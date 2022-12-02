package turbinerb

import (
	"context"
	"errors"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/cmd/meroxa/turbine/ruby/mock"
	"github.com/meroxa/cli/log"
)

func Test_Run(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tempdir := t.TempDir()

	tests := []struct {
		name    string
		cli     *turbineRbCLI
		wantErr error
	}{
		{
			name: "fail to find app",
			cli: &turbineRbCLI{
				appPath: "/tmp",
				runServer: func() turbineServer {
					m := mock.NewMockturbineServer(ctrl)
					m.EXPECT().
						Run(gomock.Any()).
						Times(1)
					m.EXPECT().
						GracefulStop().
						Times(1)
					return m
				}(),
				logger: log.NewTestLogger(),
			},
			wantErr: errors.New("exit status 1"),
		},
		{
			name: "fail to start command",
			cli: &turbineRbCLI{
				appPath: "/nonexistentdir",
				runServer: func() turbineServer {
					m := mock.NewMockturbineServer(ctrl)
					m.EXPECT().
						Run(gomock.Any()).
						Times(1)
					m.EXPECT().
						GracefulStop().
						Times(1)
					return m
				}(),
				logger: log.NewTestLogger(),
			},
			wantErr: errors.New("chdir /nonexistentdir: no such file or directory"),
		},
		{
			name: "success",
			cli: &turbineRbCLI{
				appPath: func() string {
					if err := os.WriteFile(
						path.Join(tempdir, "app.rb"),
						[]byte(`class Turbine; def self.run; puts "it ran"; end; end`),
						0644,
					); err != nil {
						t.Fatal(err)
					}

					return tempdir
				}(),
				runServer: func() turbineServer {
					m := mock.NewMockturbineServer(ctrl)
					m.EXPECT().
						Run(gomock.Any()).
						Times(1)
					m.EXPECT().
						GracefulStop().
						Times(1)
					return m
				}(),
				logger: log.NewTestLogger(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cli.Run(ctx)
			if tc.wantErr != nil && !strings.Contains(err.Error(), tc.wantErr.Error()) {
				t.Fatalf("want: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}

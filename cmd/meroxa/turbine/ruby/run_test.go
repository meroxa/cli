package turbinerb

import (
	"context"
	"errors"
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meroxa/cli/cmd/meroxa/turbine/ruby/mock"
	"github.com/meroxa/cli/log"
)

func Test_Run(t *testing.T) {
	var (
		ctx           = context.Background()
		tempdir       = t.TempDir()
		serverStarted = make(chan bool)
	)

	tests := []struct {
		name    string
		cli     func(ctrl *gomock.Controller) *turbineRbCLI
		wantErr error
	}{
		{
			name: "fail to find app",
			cli: func(ctrl *gomock.Controller) *turbineRbCLI {
				return &turbineRbCLI{
					appPath: "/tmp",
					runServer: func() turbineServer {
						m := mock.NewMockturbineServer(ctrl)
						m.EXPECT().
							Run(gomock.Any()).
							DoAndReturn(func(_ context.Context) {
								serverStarted <- true
							}).
							Times(1)
						m.EXPECT().
							GracefulStop().
							Times(1)
						return m
					}(),
					logger: log.NewTestLogger(),
				}
			},
			wantErr: errors.New("exit status 1"),
		},
		{
			name: "fail to start command",
			cli: func(ctrl *gomock.Controller) *turbineRbCLI {
				return &turbineRbCLI{
					appPath: "/nonexistentdir",
					runServer: func() turbineServer {
						m := mock.NewMockturbineServer(ctrl)
						m.EXPECT().
							Run(gomock.Any()).
							DoAndReturn(func(_ context.Context) {
								serverStarted <- true
							}).
							Times(1)
						m.EXPECT().
							GracefulStop().
							Times(1)
						return m
					}(),
					logger: log.NewTestLogger(),
				}
			},
			wantErr: errors.New("chdir /nonexistentdir: no such file or directory"),
		},
		{
			name: "success",
			cli: func(ctrl *gomock.Controller) *turbineRbCLI {
				return &turbineRbCLI{
					appPath: func() string {
						if err := os.WriteFile(
							path.Join(tempdir, "app.rb"),
							[]byte(`class TurbineRb; def self.run; puts "it ran"; end; end`),
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
							DoAndReturn(func(_ context.Context) {
								serverStarted <- true
							}).
							Times(1)
						m.EXPECT().
							GracefulStop().
							Times(1)
						return m
					}(),
					logger: log.NewTestLogger(),
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := tc.cli(ctrl)
			err := c.Run(ctx)

			time.AfterFunc(2*time.Second, func() { serverStarted <- false })
			if ok := <-serverStarted; !ok {
				t.Fatal("runserver failed to start")
			}

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, err.Error(), tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

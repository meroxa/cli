package flink

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestFlinkJobDeployAppFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "jar", required: true},
	}

	c := builder.BuildCobraCommand(&Deploy{})

	for _, f := range expectedFlags {
		cf := c.Flags().Lookup(f.name)
		if cf == nil {
			t.Fatalf("expected flag \"%s\" to be present", f.name)
		}

		if f.shorthand != cf.Shorthand {
			t.Fatalf("expected shorthand \"%s\" got \"%s\" for flag \"%s\"", f.shorthand, cf.Shorthand, f.name)
		}

		if f.required && !utils.IsFlagRequired(cf) {
			t.Fatalf("expected flag \"%s\" to be required", f.name)
		}

		if cf.Hidden != f.hidden {
			if cf.Hidden {
				t.Fatalf("expected flag \"%s\" not to be hidden", f.name)
			} else {
				t.Fatalf("expected flag \"%s\" to be hidden", f.name)
			}
		}
	}
}

//nolint:funlen // this is a test function, splitting it would duplicate code
func TestDeployFlinkJob(t *testing.T) {
	os.Setenv("UNIT_TEST", "true")
	ctx := context.Background()
	logger := log.NewTestLogger()
	name := "solid-name"
	jar := filepath.Join(os.TempDir(), "real-jar.jar")
	f, err := os.Create(jar)
	if err != nil {
		t.Fatalf(err.Error())
	}
	_, err = f.Write([]byte("oh hello"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	retries := 0
	testErr := fmt.Errorf("nope")

	server := func(status int) *httptest.Server {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			retries++
			w.WriteHeader(status)
		}))
		return server
	}
	putURL := server(http.StatusOK).URL

	tests := []struct {
		description  string
		name         string
		jar          string
		meroxaClient func(*gomock.Controller) meroxa.Client
		err          error
	}{
		{
			description: "Successfully deploy flink job",
			name:        name,
			jar:         jar,
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				sourceInput := meroxa.CreateSourceInputV2{
					Filename: filepath.Base(jar),
				}

				client.EXPECT().
					CreateSourceV2(ctx, &sourceInput).
					Return(&meroxa.Source{GetUrl: "get-url", PutUrl: putURL}, nil)

				jobInput := &meroxa.CreateFlinkJobInput{Name: name, JarURL: "get-url"}
				client.EXPECT().CreateFlinkJob(ctx, jobInput).
					Return(&meroxa.FlinkJob{}, nil)
				return client
			},
			err: nil,
		},
		{
			description: "Fail to provide name",
			jar:         jar,
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				return client
			},
			err: fmt.Errorf("the name of your Flink Job be provided as an argument"),
		},
		{
			description: "Fail to provide jar",
			name:        name,
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				return client
			},
			err: fmt.Errorf("the path to your Flink Job jar file must be provided to the --jar flag"),
		},
		{
			description: "Fail to provide file that is a jar",
			name:        name,
			jar:         "hi.jam",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				return client
			},
			err: fmt.Errorf("please provide a JAR file to the --jar flag"),
		},
		{
			description: "Fail to get source",
			name:        name,
			jar:         jar,
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				sourceInput := meroxa.CreateSourceInputV2{
					Filename: filepath.Base(jar),
				}

				client.EXPECT().
					CreateSourceV2(ctx, &sourceInput).
					Return(&meroxa.Source{GetUrl: "get-url", PutUrl: putURL}, testErr)
				return client
			},
			err: testErr,
		},
		{
			description: "Fail to create Flink Job",
			name:        name,
			jar:         jar,
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				sourceInput := meroxa.CreateSourceInputV2{
					Filename: filepath.Base(jar),
				}

				client.EXPECT().
					CreateSourceV2(ctx, &sourceInput).
					Return(&meroxa.Source{GetUrl: "get-url", PutUrl: putURL}, nil)

				jobInput := &meroxa.CreateFlinkJobInput{Name: name, JarURL: "get-url"}
				client.EXPECT().CreateFlinkJob(ctx, jobInput).
					Return(&meroxa.FlinkJob{}, testErr)
				return client
			},
			err: testErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			cfg := config.NewInMemoryConfig()
			d := &Deploy{
				client: tc.meroxaClient(ctrl),
				logger: logger,
				config: cfg,
			}
			d.args.Name = tc.name
			d.flags.Jar = tc.jar

			err := d.Execute(ctx)
			if err != nil {
				require.NotEmptyf(t, tc.err, err.Error())
				require.Equal(t, tc.err, err)
			} else {
				require.Empty(t, tc.err)
			}
		})
	}
	os.Setenv("UNIT_TEST", "")
}

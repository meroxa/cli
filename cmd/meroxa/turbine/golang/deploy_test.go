package turbinego

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client"
	mock_client "github.com/meroxa/turbine-core/pkg/client/mock"
)

func Test_NeedsToBuild(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		cli         func(*gomock.Controller) *turbineGoCLI
		wantErr     error
		needToBuild bool
	}{
		{
			name: "Has function",
			cli: func(ctrl *gomock.Controller) *turbineGoCLI {
				return &turbineGoCLI{
					bc: func() client.Client {
						m := mock_client.NewMockClient(ctrl)
						m.EXPECT().
							HasFunctions(gomock.Any(), &emptypb.Empty{}).
							Times(1).
							Return(&wrapperspb.BoolValue{
								Value: true,
							}, nil)
						return m
					}(),
					logger: log.NewTestLogger(),
				}
			},
			needToBuild: true,
		},
		{
			name: "Doesn't have function",
			cli: func(ctrl *gomock.Controller) *turbineGoCLI {
				return &turbineGoCLI{
					bc: func() client.Client {
						m := mock_client.NewMockClient(ctrl)
						m.EXPECT().
							HasFunctions(gomock.Any(), &emptypb.Empty{}).
							Times(1).
							Return(&wrapperspb.BoolValue{
								Value: false,
							}, nil)
						return m
					}(),
					logger: log.NewTestLogger(),
				}
			},
			needToBuild: false,
		},
		{
			name: "fail to get function info",
			cli: func(ctrl *gomock.Controller) *turbineGoCLI {
				return &turbineGoCLI{
					bc: func() client.Client {
						m := mock_client.NewMockClient(ctrl)
						m.EXPECT().
							HasFunctions(gomock.Any(), &emptypb.Empty{}).
							Times(1).
							Return(nil, errors.New("something went wrong"))
						return m
					}(),
					logger: log.NewTestLogger(),
				}
			},
			wantErr: errors.New("something went wrong"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := tc.cli(ctrl)
			needToBuild, err := c.NeedsToBuild(ctx)
			if tc.wantErr != nil && !strings.Contains(err.Error(), tc.wantErr.Error()) {
				t.Fatalf("want: %v, got: %v", tc.wantErr, err)
			}
			require.Equal(t, tc.needToBuild, needToBuild)
		})
	}
}

func Test_Deploy(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		cli      func(ctrl *gomock.Controller) *turbineGoCLI
		wantErr  error
		wantSpec string
	}{
		{
			name: "get spec",
			cli: func(ctrl *gomock.Controller) *turbineGoCLI {
				return &turbineGoCLI{
					bc: func() client.Client {
						m := mock_client.NewMockClient(ctrl)
						m.EXPECT().
							GetSpec(gomock.Any(), &pb.GetSpecRequest{
								Image: "image",
							}).
							Times(1).
							Return(&pb.GetSpecResponse{
								Spec: []byte("spec"),
							}, nil)
						return m
					}(),
					logger: log.NewTestLogger(),
				}
			},
			wantSpec: "spec",
		},
		{
			name: "fail to get spec",
			cli: func(ctrl *gomock.Controller) *turbineGoCLI {
				return &turbineGoCLI{
					bc: func() client.Client {
						m := mock_client.NewMockClient(ctrl)
						m.EXPECT().
							GetSpec(gomock.Any(), &pb.GetSpecRequest{
								Image: "image",
							}).
							Times(1).
							Return(nil, errors.New("something went wrong"))
						return m
					}(),
					logger: log.NewTestLogger(),
				}
			},
			wantErr: errors.New("something went wrong"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := tc.cli(ctrl)
			spec, err := c.GetDeploymentSpec(ctx, "image", "app", "git_sha", "0.2.0", "accountUUID")
			if tc.wantErr != nil && !strings.Contains(err.Error(), tc.wantErr.Error()) {
				t.Fatalf("want: %v, got: %v", tc.wantErr, err)
			}
			require.Equal(t, tc.wantSpec, spec)
		})
	}
}

func Test_GetResources(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		cli           func(ctrl *gomock.Controller) *turbineGoCLI
		wantErr       error
		wantResources []utils.ApplicationResource
	}{
		{
			name: "get spec",
			cli: func(ctrl *gomock.Controller) *turbineGoCLI {
				return &turbineGoCLI{
					bc: func() client.Client {
						m := mock_client.NewMockClient(ctrl)
						m.EXPECT().
							ListResources(gomock.Any(), &emptypb.Empty{}).
							Times(1).
							Return(&pb.ListResourcesResponse{
								Resources: []*pb.Resource{
									{
										Name: "pg",
									},
									{
										Name: "mongo",
									},
								},
							}, nil)
						return m
					}(),
					logger: log.NewTestLogger(),
				}
			},
			wantResources: []utils.ApplicationResource{
				{
					Name: "pg",
				},
				{
					Name: "mongo",
				},
			},
		},
		{
			name: "fail to list resources",
			cli: func(ctrl *gomock.Controller) *turbineGoCLI {
				return &turbineGoCLI{
					bc: func() client.Client {
						m := mock_client.NewMockClient(ctrl)
						m.EXPECT().
							ListResources(gomock.Any(), &emptypb.Empty{}).
							Times(1).
							Return(nil, errors.New("something went wrong"))
						return m
					}(),
					logger: log.NewTestLogger(),
				}
			},
			wantErr: errors.New("something went wrong"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := tc.cli(ctrl)
			spec, err := c.GetResources(ctx)
			if tc.wantErr != nil && !strings.Contains(err.Error(), tc.wantErr.Error()) {
				t.Fatalf("want: %v, got: %v", tc.wantErr, err)
			}
			require.Equal(t, tc.wantResources, spec)
		})
	}
}
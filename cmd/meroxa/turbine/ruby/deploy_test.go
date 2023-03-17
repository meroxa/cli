package turbinerb

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/mock"
	"github.com/meroxa/cli/log"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func Test_NeedsToBuild(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		cli         func(*gomock.Controller) *turbineRbCLI
		wantErr     error
		needToBuild bool
	}{
		{
			name: "Has function",
			cli: func(ctrl *gomock.Controller) *turbineRbCLI {
				return &turbineRbCLI{
					recordClient: func() recordClient {
						m := mock.NewMockTurbineServiceClient(ctrl)
						m.EXPECT().
							HasFunctions(gomock.Any(), &emptypb.Empty{}).
							Times(1).
							Return(&wrapperspb.BoolValue{
								Value: true,
							}, nil)
						return recordClient{
							TurbineServiceClient: m,
						}
					}(),
					logger: log.NewTestLogger(),
				}
			},
			needToBuild: true,
		},
		{
			name: "Doesn't have function",
			cli: func(ctrl *gomock.Controller) *turbineRbCLI {
				return &turbineRbCLI{
					recordClient: func() recordClient {
						m := mock.NewMockTurbineServiceClient(ctrl)
						m.EXPECT().
							HasFunctions(gomock.Any(), &emptypb.Empty{}).
							Times(1).
							Return(&wrapperspb.BoolValue{
								Value: false,
							}, nil)
						return recordClient{
							TurbineServiceClient: m,
						}
					}(),
					logger: log.NewTestLogger(),
				}
			},
			needToBuild: false,
		},
		{
			name: "fail to get function info",
			cli: func(ctrl *gomock.Controller) *turbineRbCLI {
				return &turbineRbCLI{
					recordClient: func() recordClient {
						m := mock.NewMockTurbineServiceClient(ctrl)
						m.EXPECT().
							HasFunctions(gomock.Any(), &emptypb.Empty{}).
							Times(1).
							Return(nil, errors.New("something went wrong"))
						return recordClient{
							TurbineServiceClient: m,
						}
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
			needToBuild, err := c.NeedsToBuild(ctx, "app")
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
		cli      func(ctrl *gomock.Controller) *turbineRbCLI
		wantErr  error
		wantSpec string
	}{
		{
			name: "get spec",
			cli: func(ctrl *gomock.Controller) *turbineRbCLI {
				return &turbineRbCLI{
					recordClient: func() recordClient {
						m := mock.NewMockTurbineServiceClient(ctrl)
						m.EXPECT().
							GetSpec(gomock.Any(), &pb.GetSpecRequest{
								Image: "image",
							}).
							Times(1).
							Return(&pb.GetSpecResponse{
								Spec: []byte("spec"),
							}, nil)
						return recordClient{
							TurbineServiceClient: m,
						}
					}(),
					logger: log.NewTestLogger(),
				}
			},
			wantSpec: "spec",
		},
		{
			name: "fail to get spec",
			cli: func(ctrl *gomock.Controller) *turbineRbCLI {
				return &turbineRbCLI{
					recordClient: func() recordClient {
						m := mock.NewMockTurbineServiceClient(ctrl)
						m.EXPECT().
							GetSpec(gomock.Any(), &pb.GetSpecRequest{
								Image: "image",
							}).
							Times(1).
							Return(nil, errors.New("something went wrong"))
						return recordClient{
							TurbineServiceClient: m,
						}
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
			spec, err := c.Deploy(ctx, "image", "app", "git_sha", "0.2.0", "accountUUID")
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
		cli           func(ctrl *gomock.Controller) *turbineRbCLI
		wantErr       error
		wantResources []utils.ApplicationResource
	}{
		{
			name: "get spec",
			cli: func(ctrl *gomock.Controller) *turbineRbCLI {
				return &turbineRbCLI{
					recordClient: func() recordClient {
						m := mock.NewMockTurbineServiceClient(ctrl)
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
						return recordClient{
							TurbineServiceClient: m,
						}
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
			cli: func(ctrl *gomock.Controller) *turbineRbCLI {
				return &turbineRbCLI{
					recordClient: func() recordClient {
						m := mock.NewMockTurbineServiceClient(ctrl)
						m.EXPECT().
							ListResources(gomock.Any(), &emptypb.Empty{}).
							Times(1).
							Return(nil, errors.New("something went wrong"))
						return recordClient{
							TurbineServiceClient: m,
						}
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
			spec, err := c.GetResources(ctx, "app")
			if tc.wantErr != nil && !strings.Contains(err.Error(), tc.wantErr.Error()) {
				t.Fatalf("want: %v, got: %v", tc.wantErr, err)
			}
			require.Equal(t, tc.wantResources, spec)
		})
	}
}

package turbine

import (
	"context"
	"errors"
	"strings"
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/golang/mock/gomock"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client"
	mock_client "github.com/meroxa/turbine-core/pkg/client/mock"
	"github.com/stretchr/testify/require"
)

func Test_GetDeploymentSpec(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		core     func(ctrl *gomock.Controller) *Core
		wantErr  error
		wantSpec string
	}{
		{
			name: "get spec",
			core: func(ctrl *gomock.Controller) *Core {
				return &Core{
					client: func() client.Client {
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
				}
			},
			wantSpec: "spec",
		},
		{
			name: "fail to get spec",
			core: func(ctrl *gomock.Controller) *Core {
				return &Core{
					client: func() client.Client {
						m := mock_client.NewMockClient(ctrl)
						m.EXPECT().
							GetSpec(gomock.Any(), &pb.GetSpecRequest{
								Image: "image",
							}).
							Times(1).
							Return(nil, errors.New("something went wrong"))
						return m
					}(),
				}
			},
			wantErr: errors.New("something went wrong"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := tc.core(ctrl)
			spec, err := c.GetDeploymentSpec(ctx, "image")
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
		core          func(ctrl *gomock.Controller) *Core
		wantErr       error
		wantResources []ApplicationResource
	}{
		{
			name: "get spec",
			core: func(ctrl *gomock.Controller) *Core {
				return &Core{
					client: func() client.Client {
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
				}
			},
			wantResources: []ApplicationResource{
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
			core: func(ctrl *gomock.Controller) *Core {
				return &Core{
					client: func() client.Client {
						m := mock_client.NewMockClient(ctrl)
						m.EXPECT().
							ListResources(gomock.Any(), &emptypb.Empty{}).
							Times(1).
							Return(nil, errors.New("something went wrong"))
						return m
					}(),
				}
			},
			wantErr: errors.New("something went wrong"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := tc.core(ctrl)
			spec, err := c.GetResources(ctx)
			if tc.wantErr != nil && !strings.Contains(err.Error(), tc.wantErr.Error()) {
				t.Fatalf("want: %v, got: %v", tc.wantErr, err)
			}
			require.Equal(t, tc.wantResources, spec)
		})
	}
}

func Test_NeedsToBuild(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		core        func(*gomock.Controller) *Core
		wantErr     error
		needToBuild bool
	}{
		{
			name: "Has function",
			core: func(ctrl *gomock.Controller) *Core {
				return &Core{
					client: func() client.Client {
						m := mock_client.NewMockClient(ctrl)
						m.EXPECT().
							HasFunctions(gomock.Any(), &emptypb.Empty{}).
							Times(1).
							Return(&wrapperspb.BoolValue{
								Value: true,
							}, nil)
						return m
					}(),
				}
			},
			needToBuild: true,
		},
		{
			name: "Doesn't have function",
			core: func(ctrl *gomock.Controller) *Core {
				return &Core{
					client: func() client.Client {
						m := mock_client.NewMockClient(ctrl)
						m.EXPECT().
							HasFunctions(gomock.Any(), &emptypb.Empty{}).
							Times(1).
							Return(&wrapperspb.BoolValue{
								Value: false,
							}, nil)
						return m
					}(),
				}
			},
			needToBuild: false,
		},
		{
			name: "fail to get function info",
			core: func(ctrl *gomock.Controller) *Core {
				return &Core{
					client: func() client.Client {
						m := mock_client.NewMockClient(ctrl)
						m.EXPECT().
							HasFunctions(gomock.Any(), &emptypb.Empty{}).
							Times(1).
							Return(nil, errors.New("something went wrong"))
						return m
					}(),
				}
			},
			wantErr: errors.New("something went wrong"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := tc.core(ctrl)
			needToBuild, err := c.NeedsToBuild(ctx)
			if tc.wantErr != nil && !strings.Contains(err.Error(), tc.wantErr.Error()) {
				t.Fatalf("want: %v, got: %v", tc.wantErr, err)
			}
			require.Equal(t, tc.needToBuild, needToBuild)
		})
	}
}

package turbine

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/turbine-core/v2/pkg/client"
	mock_client "github.com/meroxa/turbine-core/v2/pkg/client/mock"
	pb "github.com/meroxa/turbine-core/v2/proto/turbine/v2"
	"github.com/stretchr/testify/require"
)

func Test_GetDeploymentSpecV2(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		core     func(ctrl *gomock.Controller) *CoreV2
		wantErr  error
		wantSpec string
	}{
		{
			name: "get spec",
			core: func(ctrl *gomock.Controller) *CoreV2 {
				return &CoreV2{
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
			core: func(ctrl *gomock.Controller) *CoreV2 {
				return &CoreV2{
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

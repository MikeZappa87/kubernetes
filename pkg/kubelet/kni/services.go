package kni

import (
	"context"

	"github.com/MikeZappa87/kni-api/pkg/apis/runtime/beta"
	"google.golang.org/grpc"
)

type KNIService interface {
	AttachNetwork(ctx context.Context, in *beta.AttachNetworkRequest,
		opts ...grpc.CallOption) (*beta.AttachNetworkResponse, error)

	DetachNetwork(ctx context.Context, sandboxId string) error
	
	QueryPodNetwork(ctx context.Context, sandboxId string) (*beta.QueryPodNetworkResponse, error)

	SetupNodeNetwork(ctx context.Context, in *beta.SetupNodeNetworkRequest,
		 opts ...grpc.CallOption) (*beta.SetupNodeNetworkResponse, error)

	QueryNodeNetworks(ctx context.Context) (*beta.QueryNodeNetworksResponse, error)

	Up() bool
}
package testing

import (
	"context"

	"github.com/MikeZappa87/kni-api/pkg/apis/runtime/beta"
	"google.golang.org/grpc"
)

// FakeNetworkRuntimeService is a mock of KNIClient interface.
type FakeNetworkRuntimeService struct {
}

// NewMockKNIClient creates a new mock instance.
func NewNetworkRuntimeService() *FakeNetworkRuntimeService {
	return &FakeNetworkRuntimeService{}
}

func (m *FakeNetworkRuntimeService) AttachNetwork(ctx context.Context, in *beta.AttachNetworkRequest, opts ...grpc.CallOption) (*beta.AttachNetworkResponse, error) {
	return &beta.AttachNetworkResponse{}, nil
}

func (m *FakeNetworkRuntimeService) DetachNetwork(ctx context.Context, sandBoxId string) error {
	return nil
}

func (m *FakeNetworkRuntimeService) QueryNodeNetworks(ctx context.Context) (*beta.QueryNodeNetworksResponse, error) {
	return &beta.QueryNodeNetworksResponse{}, nil
}

func (m *FakeNetworkRuntimeService) QueryPodNetwork(ctx context.Context, sandboxId string) (*beta.QueryPodNetworkResponse, error) {
	return &beta.QueryPodNetworkResponse{}, nil
}

func (m *FakeNetworkRuntimeService) SetupNodeNetwork(ctx context.Context, in *beta.SetupNodeNetworkRequest, opts ...grpc.CallOption) (*beta.SetupNodeNetworkResponse, error) {
	return &beta.SetupNodeNetworkResponse{}, nil
}

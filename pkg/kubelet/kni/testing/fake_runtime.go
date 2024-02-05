package testing

import (
	"context"

	"github.com/MikeZappa87/kni-api/pkg/apis/runtime/beta"
	"google.golang.org/grpc"
)

var FakePodSandboxIPs = []string{"192.168.192.168"}

// FakeNetworkRuntimeService is a mock of KNIClient interface.
type FakeNetworkRuntimeService struct {
}

// NewMockKNIClient creates a new mock instance.
func NewNetworkRuntimeService() *FakeNetworkRuntimeService {
	return &FakeNetworkRuntimeService{}
}

func (m *FakeNetworkRuntimeService) AttachNetwork(ctx context.Context, in *beta.AttachNetworkRequest, opts ...grpc.CallOption) (*beta.AttachNetworkResponse, error) {
	return &beta.AttachNetworkResponse{
		Ipconfigs: genFakeIPConfig(),
	}, nil
}

func (m *FakeNetworkRuntimeService) DetachNetwork(ctx context.Context, sandBoxId string) error {
	return nil
}

func (m *FakeNetworkRuntimeService) QueryNodeNetworks(ctx context.Context) (*beta.QueryNodeNetworksResponse, error) {
	return &beta.QueryNodeNetworksResponse{}, nil
}

func (m *FakeNetworkRuntimeService) QueryPodNetwork(ctx context.Context, sandboxId string) (*beta.QueryPodNetworkResponse, error) {
	return &beta.QueryPodNetworkResponse{
		Ipconfigs: genFakeIPConfig(),
	}, nil
}

func (m *FakeNetworkRuntimeService) SetupNodeNetwork(ctx context.Context, in *beta.SetupNodeNetworkRequest, opts ...grpc.CallOption) (*beta.SetupNodeNetworkResponse, error) {
	return &beta.SetupNodeNetworkResponse{}, nil
}

func genFakeIPConfig() map[string] *beta.IPConfig {
	ip := make(map[string]*beta.IPConfig)

	ip["eth0"] = &beta.IPConfig{
		Ip: FakePodSandboxIPs,
	}
	return ip
}

func (m *FakeNetworkRuntimeService) Up() bool {
	return true
}

func (m *FakeNetworkRuntimeService) CreateNetwork(ctx context.Context, namespace, name string) (*beta.CreateNetworkResponse, error) {
	return &beta.CreateNetworkResponse{}, nil
}

func (m *FakeNetworkRuntimeService) DeleteNetworkById(ctx context.Context, podSandBoxID string) error {
	return nil
}

func (m *FakeNetworkRuntimeService) DeleteNetworkByPodName(ctx context.Context, name, namespace string) error {
	return nil
}
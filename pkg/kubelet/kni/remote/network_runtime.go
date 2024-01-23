package networkremote

import (
	"context"
	"os"
	"time"

	"github.com/MikeZappa87/kni-api/pkg/apis/runtime/beta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/kubelet/kni"
	"k8s.io/kubernetes/pkg/kubelet/util"
)

const (
	// connection parameters
	maxBackoffDelay      = 3 * time.Second
	baseBackoffDelay     = 100 * time.Millisecond
	minConnectionTimeout = 5 * time.Second
	maxMsgSize           = 1024 * 1024 * 16
)

type KNINetworkService struct {
	beta.KNIClient
}

// NewRemoteRuntimeService creates a new networkremote.KNIService.
func NewNetworkRuntimeService(endpoint string, connectionTimeout time.Duration) (kni.KNIService, error) {
	klog.V(3).InfoS("Connecting to runtime service", "endpoint", endpoint)
	addr, dialer, err := util.GetAddressAndDialer(endpoint)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))

	connParams := grpc.ConnectParams{
		Backoff: backoff.DefaultConfig,
	}
	connParams.MinConnectTimeout = minConnectionTimeout
	connParams.Backoff.BaseDelay = baseBackoffDelay
	connParams.Backoff.MaxDelay = maxBackoffDelay
	dialOpts = append(dialOpts,
		grpc.WithConnectParams(connParams),
	)

	conn, err := grpc.DialContext(ctx, addr, dialOpts...)
	if err != nil {
		klog.ErrorS(err, "Connect remote runtime failed", "address", addr)
		return nil, err
	}

	kni := KNINetworkService{
		beta.NewKNIClient(conn),
	}

	return &kni, nil
}

func (m *KNINetworkService) AttachNetwork(ctx context.Context, in *beta.AttachNetworkRequest, opts ...grpc.CallOption) (*beta.AttachNetworkResponse, error) {
	return m.KNIClient.AttachNetwork(ctx, in)
}

func (m *KNINetworkService) DetachNetwork(ctx context.Context, sandBoxId string) error {

	det := &beta.DetachNetworkRequest{
		Id: sandBoxId,
	}

	_, err := m.KNIClient.DetachNetwork(ctx, det)

	if err != nil {
		return err
	}

	return nil
}

func (m *KNINetworkService) QueryNodeNetworks(ctx context.Context) (*beta.QueryNodeNetworksResponse, error) {

	query := beta.QueryNodeNetworksRequest{}

	return m.KNIClient.QueryNodeNetworks(ctx, &query)
}

func (m *KNINetworkService) QueryPodNetwork(ctx context.Context, sandboxId string) (*beta.QueryPodNetworkResponse, error) {

	query := &beta.QueryPodNetworkRequest{
		Id: sandboxId,
	}

	return m.KNIClient.QueryPodNetwork(ctx, query)
}

func (m *KNINetworkService) SetupNodeNetwork(ctx context.Context, in *beta.SetupNodeNetworkRequest, opts ...grpc.CallOption) (*beta.SetupNodeNetworkResponse, error) {
	return m.KNIClient.SetupNodeNetwork(ctx, in)
}

func (m *KNINetworkService) Up() bool {
	if _, err := os.Stat("/tmp/kni.sock"); err != nil  {
		return false
	}

	return true
}
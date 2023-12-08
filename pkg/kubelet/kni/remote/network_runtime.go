package networkremote

import (
	"context"
	"net"

	"github.com/MikeZappa87/kni-server-client-example/pkg/apis/runtime/beta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type networkRuntimeService struct {
	runtimeClient beta.KNIClient
}
//protocol unix address /tmp/kni.sock
func NewNetworkRuntimeService(protocol, sockAddr string) (beta.KNIClient, error) {

	var (
		credentials = insecure.NewCredentials()
		dialer      = func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, protocol, addr)
		}
		options = []grpc.DialOption{
			grpc.WithTransportCredentials(credentials),
			grpc.WithBlock(),
			grpc.WithContextDialer(dialer),
		}
	)

	conn, err := grpc.Dial(sockAddr, options...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	return beta.NewKNIClient(conn), nil
}

func (r *networkRuntimeService) AttachNetwork(req *beta.AttachNetworkRequest) (beta.AttachNetworkResponse, error) {
	return beta.AttachNetworkResponse{}, nil
}

func (r *networkRuntimeService) DetachNetwork(req *beta.DetachNetworkRequest) (beta.DetachNetworkResponse, error) {
	return beta.DetachNetworkResponse{}, nil
}
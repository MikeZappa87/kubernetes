package networkremote

import (
	"context"
	"net"

	"github.com/MikeZappa87/kni-server-client-example/pkg/apis/runtime/beta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//protocol unix address /tmp/kni.sock
func NewNetworkRuntimeService(sockAddr string) (beta.KNIClient, error) {

	protocol := "unix"

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
		return nil, err
	}
	defer conn.Close()

	return beta.NewKNIClient(conn), nil
}

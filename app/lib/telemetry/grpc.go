// Content managed by Project Forge, see [projectforge.md] for details.
package telemetry

import (
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/encoding/gzip" // zip responses
	"google.golang.org/grpc/keepalive"
)

func NewGRPCServer(maxMsgSize int, opts ...grpc.ServerOption) *grpc.Server {
	allOpts := append([]grpc.ServerOption{
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	}, opts...)
	return grpc.NewServer(allOpts...)
}

func NewGRPCConnection(address string, maxMsgSize int, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	si := grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor())
	ui := grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor())
	ka := grpc.WithKeepaliveParams(keepalive.ClientParameters{Time: 1 * time.Minute, Timeout: 10 * time.Minute, PermitWithoutStream: true})
	insec := grpc.WithTransportCredentials(insecure.NewCredentials())
	size := grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize))
	allOpts := append([]grpc.DialOption{insec, ka, si, ui, size}, opts...)
	return grpc.Dial(address, allOpts...)
}

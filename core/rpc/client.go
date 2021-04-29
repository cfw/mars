package rpc

import (
	"context"
	"github.com/cfw/mars/core/rpc/balancer"
	"github.com/cfw/mars/core/rpc/discov/discovk8s"
	"google.golang.org/grpc"
)

func NewGrpcClient(ctx context.Context, target, balancerName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append([]grpc.DialOption{balancer.Format(balancerName)}, opts...)
	return grpc.DialContext(ctx, target, opts...)
}

func init() {
	discovk8s.RegisterResolver()
}

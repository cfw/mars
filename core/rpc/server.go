package rpc

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGrpcServer() *grpc.Server {
	grpcRecoveryOption := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			e := errors.WithStack(errors.New(fmt.Sprintf("%s", p)))
			log.Errorf("%+v", e)
			return status.Errorf(codes.Internal, "%s", p)
		}),
	}

	return grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(grpcRecoveryOption...))))
}

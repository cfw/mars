package balancer

import (
	"fmt"
	"google.golang.org/grpc"
)

func Format(name string) grpc.DialOption {
	return grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, name))
}

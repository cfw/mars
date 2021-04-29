package discovk8s

import "google.golang.org/grpc/resolver"

var disK8sBuilder builder

func RegisterResolver() {
	resolver.Register(&disK8sBuilder)
}

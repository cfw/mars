package discovk8s

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"

	"google.golang.org/grpc/resolver"
)

const subsetSize = 32

type builder struct {
	registry Registry
	once     sync.Once
}

func NewResolver() resolver.Builder {
	return &builder{}
}

type nopResolver struct {
	cc resolver.ClientConn
}

func (r *nopResolver) Close() {
}

func (r *nopResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

func (d *builder) parseTarget(target resolver.Target) (*Service, error) {
	// k8s://default/service:port
	endpoint := target.Endpoint
	ns := target.Authority
	// k8s://service.default:port/
	if endpoint == "" {
		endpoint = target.Authority
		ns = ""
	}
	s := Service{}
	if endpoint == "" {
		return nil, fmt.Errorf("target(%q) is empty", target)
	}
	var name string
	var port string
	if strings.LastIndex(endpoint, ":") < 0 {
		name = endpoint
	} else {
		var err error
		name, port, err = net.SplitHostPort(endpoint)
		if err != nil {
			return nil, fmt.Errorf("target endpoint='%s' is invalid. grpc target is %#v, err=%v", endpoint, target, err)
		}
	}

	nameSplit := strings.SplitN(name, ".", 2)
	serviceName := name
	if len(nameSplit) == 2 {
		serviceName = nameSplit[0]
		ns = nameSplit[1]
	}
	s.Name = serviceName
	s.Namespace = ns

	intPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	s.Port = int32(intPort)

	return &s, nil
}

func (d *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	d.once.Do(func() {
		client, err := NewClient()
		if err != nil {
			log.Fatalf("New k8s client error %s", err)
		}

		controller, err := NewEndpointController(client)
		if err != nil {
			log.Fatalf("New endpoint controller error %s", err)
		}

		d.registry = NewRegistry(controller)
	})

	service, err := d.parseTarget(target)
	if err != nil {
		log.Errorf("parse k8s service error: %v", err)
	}

	sub := d.registry.NewSubscriber(service)

	update := func() {
		var addrs []resolver.Address
		for _, val := range subset(sub.Values(), subsetSize) {
			addrs = append(addrs, resolver.Address{
				Addr: val,
			})
		}
		cc.UpdateState(resolver.State{
			Addresses: addrs,
		})
	}
	sub.SetUpdateFunc(update)

	update()

	return &nopResolver{cc: cc}, nil
}

func (d *builder) Scheme() string {
	return "k8s"
}

func subset(set []string, sub int) []string {
	rand.Shuffle(len(set), func(i, j int) {
		set[i], set[j] = set[j], set[i]
	})
	if len(set) <= sub {
		return set
	}

	return set[:sub]
}

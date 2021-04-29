package discovk8s

import (
	"sync"
)

type Registry interface {
	NewSubscriber(service *Service) Subscriber
}

type k8sRegistry struct {
	Registry

	controller EndpointController

	services map[string][]*ServiceInstance

	lock sync.Mutex
}

func NewRegistry(controller EndpointController) Registry {

	registry := k8sRegistry{
		controller: controller,
		services:   make(map[string][]*ServiceInstance),
	}

	return &registry
}

func (r *k8sRegistry) NewSubscriber(service *Service) Subscriber {

	serviceFullName := service.EndpointName()

	k8sSubscriber := K8sSubscriber{
		k8sRegistry: r,
		instance:    service,
	}

	r.controller.AddOnUpdateFunc(serviceFullName, func(addresses []*ServiceInstance) {
		r.lock.Lock()
		r.services[serviceFullName] = addresses
		r.lock.Unlock()

		k8sSubscriber.OnUpdate()
	})

	return &k8sSubscriber
}

func (r *k8sRegistry) GetServices(service *Service) []*ServiceInstance {
	r.lock.Lock()
	defer r.lock.Unlock()

	if value, ok := r.services[service.EndpointName()]; ok {
		var ret []*ServiceInstance

		for _, item := range value {
			if item.Port == service.Port {
				ret = append(ret, item)
			}
		}
		return ret
	} else {
		fullService, _ := r.controller.GetEndpoints(service.Name, service.Namespace)
		var ret []*ServiceInstance

		for _, item := range fullService {
			if item.Port == service.Port {
				ret = append(ret, item)
			}
		}
		return ret
	}

}

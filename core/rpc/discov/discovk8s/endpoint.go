package discovk8s

import (
	"errors"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"reflect"
	"sort"
	"sync"
)

type OnUpdateFunc func(addresses []*ServiceInstance)

type EndpointController interface {
	AddOnUpdateFunc(serviceName string, updateFunc OnUpdateFunc)
	GetEndpoints(serviceName string, namespace string) ([]*ServiceInstance, error)
}

type endpointController struct {
	endpointsInformer cache.SharedIndexInformer
	endpointsLister   listerv1.EndpointsLister
	updatefuncs       map[string][]OnUpdateFunc
	lock              sync.Mutex
}

func NewEndpointController(client *Client) (EndpointController, error) {

	informer := client.InformerFactory.Core().V1().Endpoints().Informer()
	lister := client.InformerFactory.Core().V1().Endpoints().Lister()

	c := endpointController{
		endpointsInformer: informer,
		endpointsLister:   lister,
		updatefuncs:       make(map[string][]OnUpdateFunc),
	}

	informer.AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.endpointAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.endpointUpdate,
			// Called on resource deletion.
			DeleteFunc: c.endpointDelete,
		},
	)

	client.InformerFactory.Start(client.stop)

	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(client.stop, informer.HasSynced) {
		return nil, errors.New("failed to sync k8s")
	}

	return &c, nil
}

func (e *endpointController) AddOnUpdateFunc(key string, updateFunc OnUpdateFunc) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if _, ok := e.updatefuncs[key]; !ok {
		e.updatefuncs[key] = make([]OnUpdateFunc, 1)
	}

	e.updatefuncs[key] = append(e.updatefuncs[key], updateFunc)
}

func (e *endpointController) GetEndpoints(name string, namespace string) ([]*ServiceInstance, error) {
	endpoint, err := e.endpointsLister.Endpoints(namespace).Get(name)
	if err != nil {
		log.Errorf("get endpoints error, %v", err)
		return nil, err
	}

	return e.getReadyAddress(endpoint), nil
}

func (e *endpointController) endpointAdd(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)

	epName := e.buildEndpointName(endpoints.Name, endpoints.Namespace)
	if _, ok := e.updatefuncs[epName]; !ok {
		return
	}

	e.lock.Lock()
	funcs := append([]OnUpdateFunc(nil), e.updatefuncs[epName]...)
	e.lock.Unlock()

	for _, f := range funcs {
		if f != nil {
			f(e.getReadyAddress(endpoints))
		}
	}
}

func (e *endpointController) endpointUpdate(old, new interface{}) {
	oldEndpoint := old.(*v1.Endpoints)
	newEndpoint := new.(*v1.Endpoints)

	epName := e.buildEndpointName(oldEndpoint.Name, newEndpoint.Namespace)

	if _, ok := e.updatefuncs[epName]; !ok {
		return
	}

	oldAddress := e.getReadyAddress(oldEndpoint)
	newAddress := e.getReadyAddress(newEndpoint)

	if reflect.DeepEqual(oldAddress, newAddress) {
		return
	}

	e.lock.Lock()
	var funcs []OnUpdateFunc
	funcs = append(funcs, e.updatefuncs[epName]...)
	e.lock.Unlock()

	for _, f := range funcs {
		if f != nil {
			f(newAddress)
		}
	}

}

func (e *endpointController) endpointDelete(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)

	epName := e.buildEndpointName(endpoints.Name, endpoints.Namespace)
	if _, ok := e.updatefuncs[epName]; !ok {
		return
	}
	e.lock.Lock()
	funcs := append([]OnUpdateFunc(nil), e.updatefuncs[epName]...)
	e.lock.Unlock()

	for _, f := range funcs {
		if f != nil {
			f(nil)
		}
	}
}

func (e *endpointController) getReadyAddress(endpoints *v1.Endpoints) []*ServiceInstance {
	var readyAddresses []*ServiceInstance

	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			for _, port := range subset.Ports {
				si := ServiceInstance{
					Ip:   address.IP,
					Port: port.Port,
				}
				readyAddresses = append(readyAddresses, &si)
			}
		}
	}

	sort.Sort(ServiceInstanceSlice(readyAddresses))

	return readyAddresses
}

func (e *endpointController) buildEndpointName(name string, namespace string) string {
	return name + "." + namespace
}

package discovery

import (
	"context"
	"sync"
	"time"
)

type ServiceName string
type InstanceID string

type InMemoryRegistry struct {
	sync.RWMutex
	serviceAddrs map[ServiceName]map[InstanceID]*serviceInstance
}

type serviceInstance struct {
	remoteAddr string
	lastActive time.Time
}

func NewInMemoryRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{
		serviceAddrs: map[ServiceName]map[InstanceID]*serviceInstance{},
	}
}

func (r *InMemoryRegistry) Register(ctx context.Context, serviceName string, instanceID string, remoteAddr string) error {
	r.Lock()
	defer r.Unlock()

	srvName := ServiceName(serviceName)

	if _, ok := r.serviceAddrs[srvName]; !ok {
		r.serviceAddrs[srvName] = map[InstanceID]*serviceInstance{}
	} else {
		instance := &serviceInstance{
			remoteAddr: remoteAddr,
			lastActive: time.Now(),
		}

		r.serviceAddrs[srvName][InstanceID(instanceID)] = instance
	}

	return nil
}

func (r *InMemoryRegistry) Deregister(ctx context.Context, serviceName string, instanceID string) error {
	r.Lock()
	defer r.Unlock()

	srvName := ServiceName(serviceName)
	service, ok := r.serviceAddrs[srvName]
	if !ok {
		return ErrServiceNotFound
	}

	delete(service, InstanceID(instanceID))
	return nil
}

func (r *InMemoryRegistry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	srvName := ServiceName(serviceName)
	service, ok := r.serviceAddrs[srvName]

	if !ok || len(service) == 0 {
		return nil, ErrServiceNotFound
	}

	var res []string
	for _, instance := range service {
		if instance.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, instance.remoteAddr)
	}
	return res, nil
}

func (r *InMemoryRegistry) ReportHealthyState(serviceName string, instanceID string) error {
	r.Lock()
	defer r.Unlock()

	srvName := ServiceName(serviceName)
	service, ok := r.serviceAddrs[srvName]
	if !ok {
		return ErrServiceNotFound
	}

	instance, ok := service[InstanceID(instanceID)]
	if !ok {
		return ErrInstanceNotFound
	}

	instance.lastActive = time.Now()

	return nil
}

package discovery

import (
	"context"
	"fmt"
	"net"
	"strconv"

	consul "github.com/hashicorp/consul/api"
)

type ConsulRegistry struct {
	client *consul.Client
}

func NewConsulRegistry(addr string) (*ConsulRegistry, error) {
	cconfig := consul.DefaultConfig()
	cconfig.Address = addr

	cclient, err := consul.NewClient(cconfig)

	if err != nil {
		panic(err)
	}

	return &ConsulRegistry{client: cclient}, nil
}

func (r *ConsulRegistry) Register(ctx context.Context, serviceName string, instanceID string, remoteAddr string) error {
	host, sPort, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(sPort)
	if err != nil {
		return err
	}

	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: host,
		ID:      instanceID,
		Name:    serviceName,
		Port:    port,
		Check: &consul.AgentServiceCheck{
			CheckID: instanceID,
			TTL:     "5s",
		},
	})
}

func (r *ConsulRegistry) Deregister(ctx context.Context, serviceName string, instanceID string) error {
	return r.client.Agent().ServiceDeregister(instanceID)
}

func (r *ConsulRegistry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, ErrServiceNotFound
	}

	var res []string
	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}

	return res, nil
}

func (r *ConsulRegistry) ReportHealthyState(_ string, instanceID string) error {
	return r.client.Agent().PassTTL(instanceID, "")
}

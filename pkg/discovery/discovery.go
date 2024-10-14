package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	Register(ctx context.Context, serviceName string, instanceID string, remoteAddr string) error
	Deregister(ctx context.Context, serviceName string, instanceID string) error
	ServiceAddresses(ctx context.Context, serviceName string) ([]string, error)
	// ReportHealthyState is a push mechanism for reporting
	// healthy state to the registry.
	ReportHealthyState(serviceName string, instanceID string) error
}

// map[ServiceName]map[InstanceID]*Instance
var ErrServiceNotFound = errors.New("no service found")
var ErrInstanceNotFound = errors.New("no instance found")

// GenerateInstanceID generates a pseudo-random service
// instance identifier, using a service name
// suffixed by dash and a random number.
func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}

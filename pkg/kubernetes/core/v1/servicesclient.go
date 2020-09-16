package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type ServicesClient interface {
	Services() ServicesReader
}

func NewServicesClient(resources resource.Client, filter resource.Filter) ServicesClient {
	return &servicesClient{
		Client: resources,
		filter: filter,
	}
}

type servicesClient struct {
	resource.Client
	filter resource.Filter
}

func (c *servicesClient) Services() ServicesReader {
	return NewServicesReader(c.Client, c.filter)
}

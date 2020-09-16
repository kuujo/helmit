package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type EndpointsClient interface {
	Endpoints() EndpointsReader
}

func NewEndpointsClient(resources resource.Client, filter resource.Filter) EndpointsClient {
	return &endpointsClient{
		Client: resources,
		filter: filter,
	}
}

type endpointsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *endpointsClient) Endpoints() EndpointsReader {
	return NewEndpointsReader(c.Client, c.filter)
}

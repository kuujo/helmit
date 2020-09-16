package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type EndpointsReference interface {
	Endpoints() EndpointsReader
}

func NewEndpointsReference(resources resource.Client, filter resource.Filter) EndpointsReference {
	return &endpointsReference{
		Client: resources,
		filter: filter,
	}
}

type endpointsReference struct {
	resource.Client
	filter resource.Filter
}

func (c *endpointsReference) Endpoints() EndpointsReader {
	return NewEndpointsReader(c.Client, c.filter)
}

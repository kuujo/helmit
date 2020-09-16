package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type NamespacesClient interface {
	Namespaces() NamespacesReader
}

func NewNamespacesClient(resources resource.Client, filter resource.Filter) NamespacesClient {
	return &namespacesClient{
		Client: resources,
		filter: filter,
	}
}

type namespacesClient struct {
	resource.Client
	filter resource.Filter
}

func (c *namespacesClient) Namespaces() NamespacesReader {
	return NewNamespacesReader(c.Client, c.filter)
}

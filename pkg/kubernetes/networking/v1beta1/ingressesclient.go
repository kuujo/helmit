package v1beta1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type IngressesClient interface {
	Ingresses() IngressesReader
}

func NewIngressesClient(resources resource.Client, filter resource.Filter) IngressesClient {
	return &ingressesClient{
		Client: resources,
		filter: filter,
	}
}

type ingressesClient struct {
	resource.Client
	filter resource.Filter
}

func (c *ingressesClient) Ingresses() IngressesReader {
	return NewIngressesReader(c.Client, c.filter)
}

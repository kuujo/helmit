package v1beta1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type Client interface {
	CustomResourceDefinitionsClient
}

func NewClient(resources resource.Client, filter resource.Filter) Client {
	return &client{
		Client:                          resources,
		CustomResourceDefinitionsClient: NewCustomResourceDefinitionsClient(resources, filter),
	}
}

type client struct {
	resource.Client
	CustomResourceDefinitionsClient
}

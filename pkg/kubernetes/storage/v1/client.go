package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type Client interface {
	StorageClassesClient
}

func NewClient(resources resource.Client, filter resource.Filter) Client {
	return &client{
		Client:               resources,
		StorageClassesClient: NewStorageClassesClient(resources, filter),
	}
}

type client struct {
	resource.Client
	StorageClassesClient
}

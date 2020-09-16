package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type StorageClassesClient interface {
	StorageClasses() StorageClassesReader
}

func NewStorageClassesClient(resources resource.Client, filter resource.Filter) StorageClassesClient {
	return &storageClassesClient{
		Client: resources,
		filter: filter,
	}
}

type storageClassesClient struct {
	resource.Client
	filter resource.Filter
}

func (c *storageClassesClient) StorageClasses() StorageClassesReader {
	return NewStorageClassesReader(c.Client, c.filter)
}

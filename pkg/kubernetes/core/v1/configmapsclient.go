package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type ConfigMapsClient interface {
	ConfigMaps() ConfigMapsReader
}

func NewConfigMapsClient(resources resource.Client, filter resource.Filter) ConfigMapsClient {
	return &configMapsClient{
		Client: resources,
		filter: filter,
	}
}

type configMapsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *configMapsClient) ConfigMaps() ConfigMapsReader {
	return NewConfigMapsReader(c.Client, c.filter)
}

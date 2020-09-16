package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type DeploymentsClient interface {
	Deployments() DeploymentsReader
}

func NewDeploymentsClient(resources resource.Client, filter resource.Filter) DeploymentsClient {
	return &deploymentsClient{
		Client: resources,
		filter: filter,
	}
}

type deploymentsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *deploymentsClient) Deployments() DeploymentsReader {
	return NewDeploymentsReader(c.Client, c.filter)
}

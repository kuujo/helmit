// Code generated by helmet-generate. DO NOT EDIT.

package v1beta1

import (
	"github.com/onosproject/helmet/pkg/helm/api/resource"
)

type Client interface {
	DeploymentsClient
	StatefulSetsClient
}

func NewClient(resources resource.Client, filter resource.Filter) Client {
	return &client{
		Client:             resources,
		DeploymentsClient:  NewDeploymentsClient(resources, filter),
		StatefulSetsClient: NewStatefulSetsClient(resources, filter),
	}
}

type client struct {
	resource.Client
	DeploymentsClient
	StatefulSetsClient
}

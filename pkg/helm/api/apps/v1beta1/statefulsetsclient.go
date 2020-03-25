// Code generated by helmet-generate. DO NOT EDIT.

package v1beta1

import (
	"github.com/onosproject/helmet/pkg/helm/api/resource"
)

type StatefulSetsClient interface {
	StatefulSets() StatefulSetsReader
}

func NewStatefulSetsClient(resources resource.Client, filter resource.Filter) StatefulSetsClient {
	return &statefulSetsClient{
		Client: resources,
		filter: filter,
	}
}

type statefulSetsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *statefulSetsClient) StatefulSets() StatefulSetsReader {
	return NewStatefulSetsReader(c.Client, c.filter)
}

package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
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

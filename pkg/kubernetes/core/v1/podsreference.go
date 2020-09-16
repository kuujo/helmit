package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type PodsReference interface {
	Pods() PodsReader
}

func NewPodsReference(resources resource.Client, filter resource.Filter) PodsReference {
	return &podsReference{
		Client: resources,
		filter: filter,
	}
}

type podsReference struct {
	resource.Client
	filter resource.Filter
}

func (c *podsReference) Pods() PodsReader {
	return NewPodsReader(c.Client, c.filter)
}

package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type RolesClient interface {
	Roles() RolesReader
}

func NewRolesClient(resources resource.Client, filter resource.Filter) RolesClient {
	return &rolesClient{
		Client: resources,
		filter: filter,
	}
}

type rolesClient struct {
	resource.Client
	filter resource.Filter
}

func (c *rolesClient) Roles() RolesReader {
	return NewRolesReader(c.Client, c.filter)
}

package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type RoleBindingsClient interface {
	RoleBindings() RoleBindingsReader
}

func NewRoleBindingsClient(resources resource.Client, filter resource.Filter) RoleBindingsClient {
	return &roleBindingsClient{
		Client: resources,
		filter: filter,
	}
}

type roleBindingsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *roleBindingsClient) RoleBindings() RoleBindingsReader {
	return NewRoleBindingsReader(c.Client, c.filter)
}

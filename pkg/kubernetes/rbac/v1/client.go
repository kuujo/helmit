package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type Client interface {
	ClusterRolesClient
	ClusterRoleBindingsClient
	RolesClient
	RoleBindingsClient
}

func NewClient(resources resource.Client, filter resource.Filter) Client {
	return &client{
		Client:                    resources,
		ClusterRolesClient:        NewClusterRolesClient(resources, filter),
		ClusterRoleBindingsClient: NewClusterRoleBindingsClient(resources, filter),
		RolesClient:               NewRolesClient(resources, filter),
		RoleBindingsClient:        NewRoleBindingsClient(resources, filter),
	}
}

type client struct {
	resource.Client
	ClusterRolesClient
	ClusterRoleBindingsClient
	RolesClient
	RoleBindingsClient
}

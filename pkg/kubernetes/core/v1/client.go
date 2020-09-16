package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type Client interface {
	ConfigMapsClient
	EndpointsClient
	NamespacesClient
	NodesClient
	PersistentVolumesClient
	PersistentVolumeClaimsClient
	PodsClient
	PodTemplatesClient
	SecretsClient
	ServicesClient
}

func NewClient(resources resource.Client, filter resource.Filter) Client {
	return &client{
		Client:                       resources,
		ConfigMapsClient:             NewConfigMapsClient(resources, filter),
		EndpointsClient:              NewEndpointsClient(resources, filter),
		NamespacesClient:             NewNamespacesClient(resources, filter),
		NodesClient:                  NewNodesClient(resources, filter),
		PersistentVolumesClient:      NewPersistentVolumesClient(resources, filter),
		PersistentVolumeClaimsClient: NewPersistentVolumeClaimsClient(resources, filter),
		PodsClient:                   NewPodsClient(resources, filter),
		PodTemplatesClient:           NewPodTemplatesClient(resources, filter),
		SecretsClient:                NewSecretsClient(resources, filter),
		ServicesClient:               NewServicesClient(resources, filter),
	}
}

type client struct {
	resource.Client
	ConfigMapsClient
	EndpointsClient
	NamespacesClient
	NodesClient
	PersistentVolumesClient
	PersistentVolumeClaimsClient
	PodsClient
	PodTemplatesClient
	SecretsClient
	ServicesClient
}

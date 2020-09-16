package v1beta1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type Client interface {
	PodDisruptionBudgetsClient
	PodSecurityPoliciesClient
}

func NewClient(resources resource.Client, filter resource.Filter) Client {
	return &client{
		Client:                     resources,
		PodDisruptionBudgetsClient: NewPodDisruptionBudgetsClient(resources, filter),
		PodSecurityPoliciesClient:  NewPodSecurityPoliciesClient(resources, filter),
	}
}

type client struct {
	resource.Client
	PodDisruptionBudgetsClient
	PodSecurityPoliciesClient
}

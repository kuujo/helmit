package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type Client interface {
	MutatingWebhookConfigurationsClient
	ValidatingWebhookConfigurationsClient
}

func NewClient(resources resource.Client, filter resource.Filter) Client {
	return &client{
		Client:                                resources,
		MutatingWebhookConfigurationsClient:   NewMutatingWebhookConfigurationsClient(resources, filter),
		ValidatingWebhookConfigurationsClient: NewValidatingWebhookConfigurationsClient(resources, filter),
	}
}

type client struct {
	resource.Client
	MutatingWebhookConfigurationsClient
	ValidatingWebhookConfigurationsClient
}

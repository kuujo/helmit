package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type MutatingWebhookConfigurationsClient interface {
	MutatingWebhookConfigurations() MutatingWebhookConfigurationsReader
}

func NewMutatingWebhookConfigurationsClient(resources resource.Client, filter resource.Filter) MutatingWebhookConfigurationsClient {
	return &mutatingWebhookConfigurationsClient{
		Client: resources,
		filter: filter,
	}
}

type mutatingWebhookConfigurationsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *mutatingWebhookConfigurationsClient) MutatingWebhookConfigurations() MutatingWebhookConfigurationsReader {
	return NewMutatingWebhookConfigurationsReader(c.Client, c.filter)
}

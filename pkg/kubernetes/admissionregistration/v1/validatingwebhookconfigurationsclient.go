package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type ValidatingWebhookConfigurationsClient interface {
	ValidatingWebhookConfigurations() ValidatingWebhookConfigurationsReader
}

func NewValidatingWebhookConfigurationsClient(resources resource.Client, filter resource.Filter) ValidatingWebhookConfigurationsClient {
	return &validatingWebhookConfigurationsClient{
		Client: resources,
		filter: filter,
	}
}

type validatingWebhookConfigurationsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *validatingWebhookConfigurationsClient) ValidatingWebhookConfigurations() ValidatingWebhookConfigurationsReader {
	return NewValidatingWebhookConfigurationsReader(c.Client, c.filter)
}

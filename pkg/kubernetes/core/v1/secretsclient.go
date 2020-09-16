package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type SecretsClient interface {
	Secrets() SecretsReader
}

func NewSecretsClient(resources resource.Client, filter resource.Filter) SecretsClient {
	return &secretsClient{
		Client: resources,
		filter: filter,
	}
}

type secretsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *secretsClient) Secrets() SecretsReader {
	return NewSecretsReader(c.Client, c.filter)
}

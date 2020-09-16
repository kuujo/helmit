package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type Client interface {
	JobsClient
}

func NewClient(resources resource.Client, filter resource.Filter) Client {
	return &client{
		Client:     resources,
		JobsClient: NewJobsClient(resources, filter),
	}
}

type client struct {
	resource.Client
	JobsClient
}

package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type JobsClient interface {
	Jobs() JobsReader
}

func NewJobsClient(resources resource.Client, filter resource.Filter) JobsClient {
	return &jobsClient{
		Client: resources,
		filter: filter,
	}
}

type jobsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *jobsClient) Jobs() JobsReader {
	return NewJobsReader(c.Client, c.filter)
}

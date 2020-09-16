package v2alpha1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type Client interface {
	CronJobsClient
}

func NewClient(resources resource.Client, filter resource.Filter) Client {
	return &client{
		Client:         resources,
		CronJobsClient: NewCronJobsClient(resources, filter),
	}
}

type client struct {
	resource.Client
	CronJobsClient
}

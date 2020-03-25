// Code generated by helmet-generate. DO NOT EDIT.

package v1beta1

import (
	"github.com/onosproject/helmet/pkg/helm/api/resource"
)

type CronJobsClient interface {
	CronJobs() CronJobsReader
}

func NewCronJobsClient(resources resource.Client, filter resource.Filter) CronJobsClient {
	return &cronJobsClient{
		Client: resources,
		filter: filter,
	}
}

type cronJobsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *cronJobsClient) CronJobs() CronJobsReader {
	return NewCronJobsReader(c.Client, c.filter)
}

package v1beta1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type PodDisruptionBudgetsClient interface {
	PodDisruptionBudgets() PodDisruptionBudgetsReader
}

func NewPodDisruptionBudgetsClient(resources resource.Client, filter resource.Filter) PodDisruptionBudgetsClient {
	return &podDisruptionBudgetsClient{
		Client: resources,
		filter: filter,
	}
}

type podDisruptionBudgetsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *podDisruptionBudgetsClient) PodDisruptionBudgets() PodDisruptionBudgetsReader {
	return NewPodDisruptionBudgetsReader(c.Client, c.filter)
}

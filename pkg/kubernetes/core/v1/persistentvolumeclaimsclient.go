// Code generated by helmit-generate. DO NOT EDIT.

package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type PersistentVolumeClaimsClient interface {
	PersistentVolumeClaims() PersistentVolumeClaimsReader
}

func NewPersistentVolumeClaimsClient(resources resource.Client, filter resource.Filter) PersistentVolumeClaimsClient {
	return &persistentVolumeClaimsClient{
		Client: resources,
		filter: filter,
	}
}

type persistentVolumeClaimsClient struct {
	resource.Client
	filter resource.Filter
}

func (c *persistentVolumeClaimsClient) PersistentVolumeClaims() PersistentVolumeClaimsReader {
	return NewPersistentVolumeClaimsReader(c.Client, c.filter)
}
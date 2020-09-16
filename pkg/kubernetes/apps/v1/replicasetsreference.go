package v1

import (
	corev1 "github.com/onosproject/helmit/pkg/kubernetes/core/v1"
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ReplicaSetsReference interface {
	ReplicaSets() ReplicaSetsReader
	corev1.PodsReference
}

func NewReplicaSetsReference(resources resource.Client, filter resource.Filter) ReplicaSetsReference {
	var ownerFilter resource.Filter = func(kind metav1.GroupVersionKind, meta metav1.ObjectMeta) (bool, error) {
		list, err := NewReplicaSetsReader(resources, filter).List()
		if err != nil {
			return false, err
		}
		for _, owner := range meta.OwnerReferences {
			for _, replicaSets := range list {
				if replicaSets.Object.ObjectMeta.UID == owner.UID {
					return true, nil
				}
			}
		}
		return false, nil
	}
	return &replicaSetsReference{
		Client:        resources,
		filter:        filter,
		PodsReference: corev1.NewPodsReference(resources, ownerFilter),
	}
}

type replicaSetsReference struct {
	resource.Client
	filter resource.Filter
	corev1.PodsReference
}

func (c *replicaSetsReference) ReplicaSets() ReplicaSetsReader {
	return NewReplicaSetsReader(c.Client, c.filter)
}

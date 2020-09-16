package v1

import (
	corev1 "github.com/onosproject/helmit/pkg/kubernetes/core/v1"
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var ReplicaSetKind = resource.Kind{
	Group:   "apps",
	Version: "v1",
	Kind:    "ReplicaSet",
	Scoped:  true,
}

var ReplicaSetResource = resource.Type{
	Kind: ReplicaSetKind,
	Name: "replicasets",
}

func NewReplicaSet(replicaSet *appsv1.ReplicaSet, client resource.Client) *ReplicaSet {
	return &ReplicaSet{
		Resource:      resource.NewResource(replicaSet.ObjectMeta, ReplicaSetKind, client),
		Object:        replicaSet,
		PodsReference: corev1.NewPodsReference(client, resource.NewUIDFilter(replicaSet.UID)),
	}
}

type ReplicaSet struct {
	*resource.Resource
	Object *appsv1.ReplicaSet
	corev1.PodsReference
}

func (r *ReplicaSet) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.AppsV1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, ReplicaSetKind.Scoped).
		Resource(ReplicaSetResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

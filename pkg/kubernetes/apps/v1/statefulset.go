package v1

import (
	corev1 "github.com/onosproject/helmit/pkg/kubernetes/core/v1"
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var StatefulSetKind = resource.Kind{
	Group:   "apps",
	Version: "v1",
	Kind:    "StatefulSet",
	Scoped:  true,
}

var StatefulSetResource = resource.Type{
	Kind: StatefulSetKind,
	Name: "statefulsets",
}

func NewStatefulSet(statefulSet *appsv1.StatefulSet, client resource.Client) *StatefulSet {
	return &StatefulSet{
		Resource:      resource.NewResource(statefulSet.ObjectMeta, StatefulSetKind, client),
		Object:        statefulSet,
		PodsReference: corev1.NewPodsReference(client, resource.NewUIDFilter(statefulSet.UID)),
	}
}

type StatefulSet struct {
	*resource.Resource
	Object *appsv1.StatefulSet
	corev1.PodsReference
}

func (r *StatefulSet) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.AppsV1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, StatefulSetKind.Scoped).
		Resource(StatefulSetResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

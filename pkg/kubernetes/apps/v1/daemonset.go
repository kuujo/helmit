package v1

import (
	corev1 "github.com/onosproject/helmit/pkg/kubernetes/core/v1"
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var DaemonSetKind = resource.Kind{
	Group:   "apps",
	Version: "v1",
	Kind:    "DaemonSet",
	Scoped:  true,
}

var DaemonSetResource = resource.Type{
	Kind: DaemonSetKind,
	Name: "daemonsets",
}

func NewDaemonSet(daemonSet *appsv1.DaemonSet, client resource.Client) *DaemonSet {
	return &DaemonSet{
		Resource:      resource.NewResource(daemonSet.ObjectMeta, DaemonSetKind, client),
		Object:        daemonSet,
		PodsReference: corev1.NewPodsReference(client, resource.NewUIDFilter(daemonSet.UID)),
	}
}

type DaemonSet struct {
	*resource.Resource
	Object *appsv1.DaemonSet
	corev1.PodsReference
}

func (r *DaemonSet) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.AppsV1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, DaemonSetKind.Scoped).
		Resource(DaemonSetResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

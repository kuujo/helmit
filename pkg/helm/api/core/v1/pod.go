// Code generated by helmet-generate. DO NOT EDIT.

package v1

import (
	"github.com/onosproject/helmet/pkg/helm/api/resource"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var PodKind = resource.Kind{
	Group:   "core",
	Version: "v1",
	Kind:    "Pod",
	Scoped:  true,
}

var PodResource = resource.Type{
	Kind: PodKind,
	Name: "pods",
}

func NewPod(pod *corev1.Pod, client resource.Client) *Pod {
	return &Pod{
		Resource: resource.NewResource(pod.ObjectMeta, PodKind, client),
		Object:   pod,
	}
}

type Pod struct {
	*resource.Resource
	Object *corev1.Pod
}

func (r *Pod) Delete() error {
	return r.Clientset().
		CoreV1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, PodKind.Scoped).
		Resource(PodResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

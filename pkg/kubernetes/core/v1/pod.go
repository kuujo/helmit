package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var PodKind = resource.Kind{
	Group:   "",
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
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.CoreV1().
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

package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var NamespaceKind = resource.Kind{
	Group:   "",
	Version: "v1",
	Kind:    "Namespace",
	Scoped:  false,
}

var NamespaceResource = resource.Type{
	Kind: NamespaceKind,
	Name: "namespaces",
}

func NewNamespace(namespace *corev1.Namespace, client resource.Client) *Namespace {
	return &Namespace{
		Resource: resource.NewResource(namespace.ObjectMeta, NamespaceKind, client),
		Object:   namespace,
	}
}

type Namespace struct {
	*resource.Resource
	Object *corev1.Namespace
}

func (r *Namespace) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.CoreV1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, NamespaceKind.Scoped).
		Resource(NamespaceResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

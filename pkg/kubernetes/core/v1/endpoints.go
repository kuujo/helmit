package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var EndpointsKind = resource.Kind{
	Group:   "",
	Version: "v1",
	Kind:    "Endpoints",
	Scoped:  true,
}

var EndpointsResource = resource.Type{
	Kind: EndpointsKind,
	Name: "endpoints",
}

func NewEndpoints(endpoints *corev1.Endpoints, client resource.Client) *Endpoints {
	return &Endpoints{
		Resource: resource.NewResource(endpoints.ObjectMeta, EndpointsKind, client),
		Object:   endpoints,
	}
}

type Endpoints struct {
	*resource.Resource
	Object *corev1.Endpoints
}

func (r *Endpoints) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.CoreV1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, EndpointsKind.Scoped).
		Resource(EndpointsResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

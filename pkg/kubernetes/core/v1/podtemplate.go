package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var PodTemplateKind = resource.Kind{
	Group:   "",
	Version: "v1",
	Kind:    "PodTemplate",
	Scoped:  true,
}

var PodTemplateResource = resource.Type{
	Kind: PodTemplateKind,
	Name: "podtemplates",
}

func NewPodTemplate(podTemplate *corev1.PodTemplate, client resource.Client) *PodTemplate {
	return &PodTemplate{
		Resource: resource.NewResource(podTemplate.ObjectMeta, PodTemplateKind, client),
		Object:   podTemplate,
	}
}

type PodTemplate struct {
	*resource.Resource
	Object *corev1.PodTemplate
}

func (r *PodTemplate) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.CoreV1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, PodTemplateKind.Scoped).
		Resource(PodTemplateResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

type PodTemplatesReader interface {
	Get(name string) (*PodTemplate, error)
	List() ([]*PodTemplate, error)
}

func NewPodTemplatesReader(client resource.Client, filter resource.Filter) PodTemplatesReader {
	return &podTemplatesReader{
		Client: client,
		filter: filter,
	}
}

type podTemplatesReader struct {
	resource.Client
	filter resource.Filter
}

func (c *podTemplatesReader) Get(name string) (*PodTemplate, error) {
	podTemplate := &corev1.PodTemplate{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.CoreV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), PodTemplateKind.Scoped).
		Resource(PodTemplateResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(podTemplate)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   PodTemplateKind.Group,
			Version: PodTemplateKind.Version,
			Kind:    PodTemplateKind.Kind,
		}, podTemplate.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    PodTemplateKind.Group,
				Resource: PodTemplateResource.Name,
			}, name)
		}
	}
	return NewPodTemplate(podTemplate, c.Client), nil
}

func (c *podTemplatesReader) List() ([]*PodTemplate, error) {
	list := &corev1.PodTemplateList{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.CoreV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), PodTemplateKind.Scoped).
		Resource(PodTemplateResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*PodTemplate, 0, len(list.Items))
	for _, podTemplate := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   PodTemplateKind.Group,
			Version: PodTemplateKind.Version,
			Kind:    PodTemplateKind.Kind,
		}, podTemplate.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := podTemplate
			results = append(results, NewPodTemplate(&copy, c.Client))
		}
	}
	return results, nil
}

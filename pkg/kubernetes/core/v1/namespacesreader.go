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

type NamespacesReader interface {
	Get(name string) (*Namespace, error)
	List() ([]*Namespace, error)
}

func NewNamespacesReader(client resource.Client, filter resource.Filter) NamespacesReader {
	return &namespacesReader{
		Client: client,
		filter: filter,
	}
}

type namespacesReader struct {
	resource.Client
	filter resource.Filter
}

func (c *namespacesReader) Get(name string) (*Namespace, error) {
	namespace := &corev1.Namespace{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.CoreV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), NamespaceKind.Scoped).
		Resource(NamespaceResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(namespace)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   NamespaceKind.Group,
			Version: NamespaceKind.Version,
			Kind:    NamespaceKind.Kind,
		}, namespace.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    NamespaceKind.Group,
				Resource: NamespaceResource.Name,
			}, name)
		}
	}
	return NewNamespace(namespace, c.Client), nil
}

func (c *namespacesReader) List() ([]*Namespace, error) {
	list := &corev1.NamespaceList{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.CoreV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), NamespaceKind.Scoped).
		Resource(NamespaceResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*Namespace, 0, len(list.Items))
	for _, namespace := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   NamespaceKind.Group,
			Version: NamespaceKind.Version,
			Kind:    NamespaceKind.Kind,
		}, namespace.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := namespace
			results = append(results, NewNamespace(&copy, c.Client))
		}
	}
	return results, nil
}

package v1beta1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

type StatefulSetsReader interface {
	Get(name string) (*StatefulSet, error)
	List() ([]*StatefulSet, error)
}

func NewStatefulSetsReader(client resource.Client, filter resource.Filter) StatefulSetsReader {
	return &statefulSetsReader{
		Client: client,
		filter: filter,
	}
}

type statefulSetsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *statefulSetsReader) Get(name string) (*StatefulSet, error) {
	statefulSet := &appsv1beta1.StatefulSet{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.AppsV1beta1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), StatefulSetKind.Scoped).
		Resource(StatefulSetResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(statefulSet)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   StatefulSetKind.Group,
			Version: StatefulSetKind.Version,
			Kind:    StatefulSetKind.Kind,
		}, statefulSet.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    StatefulSetKind.Group,
				Resource: StatefulSetResource.Name,
			}, name)
		}
	}
	return NewStatefulSet(statefulSet, c.Client), nil
}

func (c *statefulSetsReader) List() ([]*StatefulSet, error) {
	list := &appsv1beta1.StatefulSetList{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.AppsV1beta1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), StatefulSetKind.Scoped).
		Resource(StatefulSetResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*StatefulSet, 0, len(list.Items))
	for _, statefulSet := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   StatefulSetKind.Group,
			Version: StatefulSetKind.Version,
			Kind:    StatefulSetKind.Kind,
		}, statefulSet.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := statefulSet
			results = append(results, NewStatefulSet(&copy, c.Client))
		}
	}
	return results, nil
}

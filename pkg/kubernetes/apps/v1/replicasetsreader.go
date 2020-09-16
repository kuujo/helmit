package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

type ReplicaSetsReader interface {
	Get(name string) (*ReplicaSet, error)
	List() ([]*ReplicaSet, error)
}

func NewReplicaSetsReader(client resource.Client, filter resource.Filter) ReplicaSetsReader {
	return &replicaSetsReader{
		Client: client,
		filter: filter,
	}
}

type replicaSetsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *replicaSetsReader) Get(name string) (*ReplicaSet, error) {
	replicaSet := &appsv1.ReplicaSet{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.AppsV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), ReplicaSetKind.Scoped).
		Resource(ReplicaSetResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(replicaSet)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   ReplicaSetKind.Group,
			Version: ReplicaSetKind.Version,
			Kind:    ReplicaSetKind.Kind,
		}, replicaSet.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    ReplicaSetKind.Group,
				Resource: ReplicaSetResource.Name,
			}, name)
		}
	}
	return NewReplicaSet(replicaSet, c.Client), nil
}

func (c *replicaSetsReader) List() ([]*ReplicaSet, error) {
	list := &appsv1.ReplicaSetList{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.AppsV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), ReplicaSetKind.Scoped).
		Resource(ReplicaSetResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*ReplicaSet, 0, len(list.Items))
	for _, replicaSet := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   ReplicaSetKind.Group,
			Version: ReplicaSetKind.Version,
			Kind:    ReplicaSetKind.Kind,
		}, replicaSet.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := replicaSet
			results = append(results, NewReplicaSet(&copy, c.Client))
		}
	}
	return results, nil
}

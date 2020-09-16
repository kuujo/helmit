package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

type ClusterRoleBindingsReader interface {
	Get(name string) (*ClusterRoleBinding, error)
	List() ([]*ClusterRoleBinding, error)
}

func NewClusterRoleBindingsReader(client resource.Client, filter resource.Filter) ClusterRoleBindingsReader {
	return &clusterRoleBindingsReader{
		Client: client,
		filter: filter,
	}
}

type clusterRoleBindingsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *clusterRoleBindingsReader) Get(name string) (*ClusterRoleBinding, error) {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.RbacV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), ClusterRoleBindingKind.Scoped).
		Resource(ClusterRoleBindingResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(clusterRoleBinding)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   ClusterRoleBindingKind.Group,
			Version: ClusterRoleBindingKind.Version,
			Kind:    ClusterRoleBindingKind.Kind,
		}, clusterRoleBinding.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    ClusterRoleBindingKind.Group,
				Resource: ClusterRoleBindingResource.Name,
			}, name)
		}
	}
	return NewClusterRoleBinding(clusterRoleBinding, c.Client), nil
}

func (c *clusterRoleBindingsReader) List() ([]*ClusterRoleBinding, error) {
	list := &rbacv1.ClusterRoleBindingList{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.RbacV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), ClusterRoleBindingKind.Scoped).
		Resource(ClusterRoleBindingResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*ClusterRoleBinding, 0, len(list.Items))
	for _, clusterRoleBinding := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   ClusterRoleBindingKind.Group,
			Version: ClusterRoleBindingKind.Version,
			Kind:    ClusterRoleBindingKind.Kind,
		}, clusterRoleBinding.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := clusterRoleBinding
			results = append(results, NewClusterRoleBinding(&copy, c.Client))
		}
	}
	return results, nil
}

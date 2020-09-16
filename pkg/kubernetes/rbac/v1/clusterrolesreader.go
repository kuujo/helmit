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

type ClusterRolesReader interface {
	Get(name string) (*ClusterRole, error)
	List() ([]*ClusterRole, error)
}

func NewClusterRolesReader(client resource.Client, filter resource.Filter) ClusterRolesReader {
	return &clusterRolesReader{
		Client: client,
		filter: filter,
	}
}

type clusterRolesReader struct {
	resource.Client
	filter resource.Filter
}

func (c *clusterRolesReader) Get(name string) (*ClusterRole, error) {
	clusterRole := &rbacv1.ClusterRole{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.RbacV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), ClusterRoleKind.Scoped).
		Resource(ClusterRoleResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(clusterRole)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   ClusterRoleKind.Group,
			Version: ClusterRoleKind.Version,
			Kind:    ClusterRoleKind.Kind,
		}, clusterRole.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    ClusterRoleKind.Group,
				Resource: ClusterRoleResource.Name,
			}, name)
		}
	}
	return NewClusterRole(clusterRole, c.Client), nil
}

func (c *clusterRolesReader) List() ([]*ClusterRole, error) {
	list := &rbacv1.ClusterRoleList{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.RbacV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), ClusterRoleKind.Scoped).
		Resource(ClusterRoleResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*ClusterRole, 0, len(list.Items))
	for _, clusterRole := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   ClusterRoleKind.Group,
			Version: ClusterRoleKind.Version,
			Kind:    ClusterRoleKind.Kind,
		}, clusterRole.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := clusterRole
			results = append(results, NewClusterRole(&copy, c.Client))
		}
	}
	return results, nil
}

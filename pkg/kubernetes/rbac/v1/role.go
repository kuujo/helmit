package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var RoleKind = resource.Kind{
	Group:   "rbac.authorization.k8s.io",
	Version: "v1",
	Kind:    "Role",
	Scoped:  true,
}

var RoleResource = resource.Type{
	Kind: RoleKind,
	Name: "roles",
}

func NewRole(role *rbacv1.Role, client resource.Client) *Role {
	return &Role{
		Resource: resource.NewResource(role.ObjectMeta, RoleKind, client),
		Object:   role,
	}
}

type Role struct {
	*resource.Resource
	Object *rbacv1.Role
}

func (r *Role) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.RbacV1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, RoleKind.Scoped).
		Resource(RoleResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

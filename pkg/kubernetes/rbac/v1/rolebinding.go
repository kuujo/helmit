package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var RoleBindingKind = resource.Kind{
	Group:   "rbac.authorization.k8s.io",
	Version: "v1",
	Kind:    "RoleBinding",
	Scoped:  true,
}

var RoleBindingResource = resource.Type{
	Kind: RoleBindingKind,
	Name: "rolebindings",
}

func NewRoleBinding(roleBinding *rbacv1.RoleBinding, client resource.Client) *RoleBinding {
	return &RoleBinding{
		Resource: resource.NewResource(roleBinding.ObjectMeta, RoleBindingKind, client),
		Object:   roleBinding,
	}
}

type RoleBinding struct {
	*resource.Resource
	Object *rbacv1.RoleBinding
}

func (r *RoleBinding) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.RbacV1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, RoleBindingKind.Scoped).
		Resource(RoleBindingResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

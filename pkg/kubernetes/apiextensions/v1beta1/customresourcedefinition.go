package v1beta1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	clientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var CustomResourceDefinitionKind = resource.Kind{
	Group:   "apiextensions.k8s.io",
	Version: "v1beta1",
	Kind:    "CustomResourceDefinition",
	Scoped:  false,
}

var CustomResourceDefinitionResource = resource.Type{
	Kind: CustomResourceDefinitionKind,
	Name: "customresourcedefinitions",
}

func NewCustomResourceDefinition(customResourceDefinition *apiextensionsv1beta1.CustomResourceDefinition, client resource.Client) *CustomResourceDefinition {
	return &CustomResourceDefinition{
		Resource: resource.NewResource(customResourceDefinition.ObjectMeta, CustomResourceDefinitionKind, client),
		Object:   customResourceDefinition,
	}
}

type CustomResourceDefinition struct {
	*resource.Resource
	Object *apiextensionsv1beta1.CustomResourceDefinition
}

func (r *CustomResourceDefinition) Delete() error {
	client, err := clientset.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.ApiextensionsV1beta1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, CustomResourceDefinitionKind.Scoped).
		Resource(CustomResourceDefinitionResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

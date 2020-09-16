package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	clientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"
)

type CustomResourceDefinitionsReader interface {
	Get(name string) (*CustomResourceDefinition, error)
	List() ([]*CustomResourceDefinition, error)
}

func NewCustomResourceDefinitionsReader(client resource.Client, filter resource.Filter) CustomResourceDefinitionsReader {
	return &customResourceDefinitionsReader{
		Client: client,
		filter: filter,
	}
}

type customResourceDefinitionsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *customResourceDefinitionsReader) Get(name string) (*CustomResourceDefinition, error) {
	customResourceDefinition := &apiextensionsv1.CustomResourceDefinition{}
	client, err := clientset.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.ApiextensionsV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), CustomResourceDefinitionKind.Scoped).
		Resource(CustomResourceDefinitionResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(customResourceDefinition)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   CustomResourceDefinitionKind.Group,
			Version: CustomResourceDefinitionKind.Version,
			Kind:    CustomResourceDefinitionKind.Kind,
		}, customResourceDefinition.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    CustomResourceDefinitionKind.Group,
				Resource: CustomResourceDefinitionResource.Name,
			}, name)
		}
	}
	return NewCustomResourceDefinition(customResourceDefinition, c.Client), nil
}

func (c *customResourceDefinitionsReader) List() ([]*CustomResourceDefinition, error) {
	list := &apiextensionsv1.CustomResourceDefinitionList{}
	client, err := clientset.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.ApiextensionsV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), CustomResourceDefinitionKind.Scoped).
		Resource(CustomResourceDefinitionResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*CustomResourceDefinition, 0, len(list.Items))
	for _, customResourceDefinition := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   CustomResourceDefinitionKind.Group,
			Version: CustomResourceDefinitionKind.Version,
			Kind:    CustomResourceDefinitionKind.Kind,
		}, customResourceDefinition.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := customResourceDefinition
			results = append(results, NewCustomResourceDefinition(&copy, c.Client))
		}
	}
	return results, nil
}

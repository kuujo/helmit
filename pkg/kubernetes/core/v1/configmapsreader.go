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

type ConfigMapsReader interface {
	Get(name string) (*ConfigMap, error)
	List() ([]*ConfigMap, error)
}

func NewConfigMapsReader(client resource.Client, filter resource.Filter) ConfigMapsReader {
	return &configMapsReader{
		Client: client,
		filter: filter,
	}
}

type configMapsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *configMapsReader) Get(name string) (*ConfigMap, error) {
	configMap := &corev1.ConfigMap{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.CoreV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), ConfigMapKind.Scoped).
		Resource(ConfigMapResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(configMap)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   ConfigMapKind.Group,
			Version: ConfigMapKind.Version,
			Kind:    ConfigMapKind.Kind,
		}, configMap.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    ConfigMapKind.Group,
				Resource: ConfigMapResource.Name,
			}, name)
		}
	}
	return NewConfigMap(configMap, c.Client), nil
}

func (c *configMapsReader) List() ([]*ConfigMap, error) {
	list := &corev1.ConfigMapList{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.CoreV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), ConfigMapKind.Scoped).
		Resource(ConfigMapResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*ConfigMap, 0, len(list.Items))
	for _, configMap := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   ConfigMapKind.Group,
			Version: ConfigMapKind.Version,
			Kind:    ConfigMapKind.Kind,
		}, configMap.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := configMap
			results = append(results, NewConfigMap(&copy, c.Client))
		}
	}
	return results, nil
}

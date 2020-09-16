package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

type MutatingWebhookConfigurationsReader interface {
	Get(name string) (*MutatingWebhookConfiguration, error)
	List() ([]*MutatingWebhookConfiguration, error)
}

func NewMutatingWebhookConfigurationsReader(client resource.Client, filter resource.Filter) MutatingWebhookConfigurationsReader {
	return &mutatingWebhookConfigurationsReader{
		Client: client,
		filter: filter,
	}
}

type mutatingWebhookConfigurationsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *mutatingWebhookConfigurationsReader) Get(name string) (*MutatingWebhookConfiguration, error) {
	mutatingWebhookConfiguration := &admissionregistrationv1.MutatingWebhookConfiguration{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.AdmissionregistrationV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), MutatingWebhookConfigurationKind.Scoped).
		Resource(MutatingWebhookConfigurationResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(mutatingWebhookConfiguration)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   MutatingWebhookConfigurationKind.Group,
			Version: MutatingWebhookConfigurationKind.Version,
			Kind:    MutatingWebhookConfigurationKind.Kind,
		}, mutatingWebhookConfiguration.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    MutatingWebhookConfigurationKind.Group,
				Resource: MutatingWebhookConfigurationResource.Name,
			}, name)
		}
	}
	return NewMutatingWebhookConfiguration(mutatingWebhookConfiguration, c.Client), nil
}

func (c *mutatingWebhookConfigurationsReader) List() ([]*MutatingWebhookConfiguration, error) {
	list := &admissionregistrationv1.MutatingWebhookConfigurationList{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.AdmissionregistrationV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), MutatingWebhookConfigurationKind.Scoped).
		Resource(MutatingWebhookConfigurationResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*MutatingWebhookConfiguration, 0, len(list.Items))
	for _, mutatingWebhookConfiguration := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   MutatingWebhookConfigurationKind.Group,
			Version: MutatingWebhookConfigurationKind.Version,
			Kind:    MutatingWebhookConfigurationKind.Kind,
		}, mutatingWebhookConfiguration.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := mutatingWebhookConfiguration
			results = append(results, NewMutatingWebhookConfiguration(&copy, c.Client))
		}
	}
	return results, nil
}

package v1beta1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

type PodDisruptionBudgetsReader interface {
	Get(name string) (*PodDisruptionBudget, error)
	List() ([]*PodDisruptionBudget, error)
}

func NewPodDisruptionBudgetsReader(client resource.Client, filter resource.Filter) PodDisruptionBudgetsReader {
	return &podDisruptionBudgetsReader{
		Client: client,
		filter: filter,
	}
}

type podDisruptionBudgetsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *podDisruptionBudgetsReader) Get(name string) (*PodDisruptionBudget, error) {
	podDisruptionBudget := &policyv1beta1.PodDisruptionBudget{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.PolicyV1beta1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), PodDisruptionBudgetKind.Scoped).
		Resource(PodDisruptionBudgetResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(podDisruptionBudget)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   PodDisruptionBudgetKind.Group,
			Version: PodDisruptionBudgetKind.Version,
			Kind:    PodDisruptionBudgetKind.Kind,
		}, podDisruptionBudget.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    PodDisruptionBudgetKind.Group,
				Resource: PodDisruptionBudgetResource.Name,
			}, name)
		}
	}
	return NewPodDisruptionBudget(podDisruptionBudget, c.Client), nil
}

func (c *podDisruptionBudgetsReader) List() ([]*PodDisruptionBudget, error) {
	list := &policyv1beta1.PodDisruptionBudgetList{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.PolicyV1beta1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), PodDisruptionBudgetKind.Scoped).
		Resource(PodDisruptionBudgetResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*PodDisruptionBudget, 0, len(list.Items))
	for _, podDisruptionBudget := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   PodDisruptionBudgetKind.Group,
			Version: PodDisruptionBudgetKind.Version,
			Kind:    PodDisruptionBudgetKind.Kind,
		}, podDisruptionBudget.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := podDisruptionBudget
			results = append(results, NewPodDisruptionBudget(&copy, c.Client))
		}
	}
	return results, nil
}

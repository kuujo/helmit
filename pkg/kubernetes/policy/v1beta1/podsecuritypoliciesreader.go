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

type PodSecurityPoliciesReader interface {
	Get(name string) (*PodSecurityPolicy, error)
	List() ([]*PodSecurityPolicy, error)
}

func NewPodSecurityPoliciesReader(client resource.Client, filter resource.Filter) PodSecurityPoliciesReader {
	return &podSecurityPoliciesReader{
		Client: client,
		filter: filter,
	}
}

type podSecurityPoliciesReader struct {
	resource.Client
	filter resource.Filter
}

func (c *podSecurityPoliciesReader) Get(name string) (*PodSecurityPolicy, error) {
	podSecurityPolicy := &policyv1beta1.PodSecurityPolicy{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.PolicyV1beta1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), PodSecurityPolicyKind.Scoped).
		Resource(PodSecurityPolicyResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(podSecurityPolicy)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   PodSecurityPolicyKind.Group,
			Version: PodSecurityPolicyKind.Version,
			Kind:    PodSecurityPolicyKind.Kind,
		}, podSecurityPolicy.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    PodSecurityPolicyKind.Group,
				Resource: PodSecurityPolicyResource.Name,
			}, name)
		}
	}
	return NewPodSecurityPolicy(podSecurityPolicy, c.Client), nil
}

func (c *podSecurityPoliciesReader) List() ([]*PodSecurityPolicy, error) {
	list := &policyv1beta1.PodSecurityPolicyList{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.PolicyV1beta1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), PodSecurityPolicyKind.Scoped).
		Resource(PodSecurityPolicyResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*PodSecurityPolicy, 0, len(list.Items))
	for _, podSecurityPolicy := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   PodSecurityPolicyKind.Group,
			Version: PodSecurityPolicyKind.Version,
			Kind:    PodSecurityPolicyKind.Kind,
		}, podSecurityPolicy.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := podSecurityPolicy
			results = append(results, NewPodSecurityPolicy(&copy, c.Client))
		}
	}
	return results, nil
}

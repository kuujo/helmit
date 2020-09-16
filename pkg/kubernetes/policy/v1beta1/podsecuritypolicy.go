package v1beta1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var PodSecurityPolicyKind = resource.Kind{
	Group:   "policy",
	Version: "v1beta1",
	Kind:    "PodSecurityPolicy",
	Scoped:  true,
}

var PodSecurityPolicyResource = resource.Type{
	Kind: PodSecurityPolicyKind,
	Name: "podsecuritypolicies",
}

func NewPodSecurityPolicy(podSecurityPolicy *policyv1beta1.PodSecurityPolicy, client resource.Client) *PodSecurityPolicy {
	return &PodSecurityPolicy{
		Resource: resource.NewResource(podSecurityPolicy.ObjectMeta, PodSecurityPolicyKind, client),
		Object:   podSecurityPolicy,
	}
}

type PodSecurityPolicy struct {
	*resource.Resource
	Object *policyv1beta1.PodSecurityPolicy
}

func (r *PodSecurityPolicy) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.PolicyV1beta1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, PodSecurityPolicyKind.Scoped).
		Resource(PodSecurityPolicyResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

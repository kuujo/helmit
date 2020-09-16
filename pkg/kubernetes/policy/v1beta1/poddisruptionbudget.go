package v1beta1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var PodDisruptionBudgetKind = resource.Kind{
	Group:   "policy",
	Version: "v1beta1",
	Kind:    "PodDisruptionBudget",
	Scoped:  true,
}

var PodDisruptionBudgetResource = resource.Type{
	Kind: PodDisruptionBudgetKind,
	Name: "poddisruptionbudgets",
}

func NewPodDisruptionBudget(podDisruptionBudget *policyv1beta1.PodDisruptionBudget, client resource.Client) *PodDisruptionBudget {
	return &PodDisruptionBudget{
		Resource: resource.NewResource(podDisruptionBudget.ObjectMeta, PodDisruptionBudgetKind, client),
		Object:   podDisruptionBudget,
	}
}

type PodDisruptionBudget struct {
	*resource.Resource
	Object *policyv1beta1.PodDisruptionBudget
}

func (r *PodDisruptionBudget) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.PolicyV1beta1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, PodDisruptionBudgetKind.Scoped).
		Resource(PodDisruptionBudgetResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

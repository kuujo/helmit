package v1beta1

import (
	appsv1 "github.com/onosproject/helmit/pkg/kubernetes/apps/v1"
	corev1 "github.com/onosproject/helmit/pkg/kubernetes/core/v1"
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var DeploymentKind = resource.Kind{
	Group:   "apps",
	Version: "v1beta1",
	Kind:    "Deployment",
	Scoped:  true,
}

var DeploymentResource = resource.Type{
	Kind: DeploymentKind,
	Name: "deployments",
}

func NewDeployment(deployment *appsv1beta1.Deployment, client resource.Client) *Deployment {
	return &Deployment{
		Resource:             resource.NewResource(deployment.ObjectMeta, DeploymentKind, client),
		Object:               deployment,
		ReplicaSetsReference: appsv1.NewReplicaSetsReference(client, resource.NewUIDFilter(deployment.UID)),
		PodsReference:        corev1.NewPodsReference(client, resource.NewUIDFilter(deployment.UID)),
	}
}

type Deployment struct {
	*resource.Resource
	Object *appsv1beta1.Deployment
	appsv1.ReplicaSetsReference
	corev1.PodsReference
}

func (r *Deployment) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.AppsV1beta1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, DeploymentKind.Scoped).
		Resource(DeploymentResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var DeploymentKind = resource.Kind{
	Group:   "apps",
	Version: "v1",
	Kind:    "Deployment",
	Scoped:  true,
}

var DeploymentResource = resource.Type{
	Kind: DeploymentKind,
	Name: "deployments",
}

func NewDeployment(deployment *appsv1.Deployment, client resource.Client) *Deployment {
	return &Deployment{
		Resource:             resource.NewResource(deployment.ObjectMeta, DeploymentKind, client),
		Object:               deployment,
		ReplicaSetsReference: NewReplicaSetsReference(client, resource.NewUIDFilter(deployment.UID)),
	}
}

type Deployment struct {
	*resource.Resource
	Object *appsv1.Deployment
	ReplicaSetsReference
}

func (r *Deployment) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.AppsV1().
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

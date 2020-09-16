package v1beta1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

type DeploymentsReader interface {
	Get(name string) (*Deployment, error)
	List() ([]*Deployment, error)
}

func NewDeploymentsReader(client resource.Client, filter resource.Filter) DeploymentsReader {
	return &deploymentsReader{
		Client: client,
		filter: filter,
	}
}

type deploymentsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *deploymentsReader) Get(name string) (*Deployment, error) {
	deployment := &appsv1beta1.Deployment{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.AppsV1beta1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), DeploymentKind.Scoped).
		Resource(DeploymentResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(deployment)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   DeploymentKind.Group,
			Version: DeploymentKind.Version,
			Kind:    DeploymentKind.Kind,
		}, deployment.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    DeploymentKind.Group,
				Resource: DeploymentResource.Name,
			}, name)
		}
	}
	return NewDeployment(deployment, c.Client), nil
}

func (c *deploymentsReader) List() ([]*Deployment, error) {
	list := &appsv1beta1.DeploymentList{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.AppsV1beta1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), DeploymentKind.Scoped).
		Resource(DeploymentResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*Deployment, 0, len(list.Items))
	for _, deployment := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   DeploymentKind.Group,
			Version: DeploymentKind.Version,
			Kind:    DeploymentKind.Kind,
		}, deployment.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := deployment
			results = append(results, NewDeployment(&copy, c.Client))
		}
	}
	return results, nil
}

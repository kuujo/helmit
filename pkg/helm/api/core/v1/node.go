// Code generated by helmet-generate. DO NOT EDIT.

package v1

import (
	"github.com/onosproject/helmet/pkg/helm/api/resource"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var NodeKind = resource.Kind{
	Group:   "core",
	Version: "v1",
	Kind:    "Node",
	Scoped:  true,
}

var NodeResource = resource.Type{
	Kind: NodeKind,
	Name: "nodes",
}

func NewNode(node *corev1.Node, client resource.Client) *Node {
	return &Node{
		Resource: resource.NewResource(node.ObjectMeta, NodeKind, client),
		Object:   node,
	}
}

type Node struct {
	*resource.Resource
	Object *corev1.Node
}

func (r *Node) Delete() error {
	return r.Clientset().
		CoreV1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, NodeKind.Scoped).
		Resource(NodeResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

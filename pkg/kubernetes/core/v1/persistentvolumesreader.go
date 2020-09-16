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

type PersistentVolumesReader interface {
	Get(name string) (*PersistentVolume, error)
	List() ([]*PersistentVolume, error)
}

func NewPersistentVolumesReader(client resource.Client, filter resource.Filter) PersistentVolumesReader {
	return &persistentVolumesReader{
		Client: client,
		filter: filter,
	}
}

type persistentVolumesReader struct {
	resource.Client
	filter resource.Filter
}

func (c *persistentVolumesReader) Get(name string) (*PersistentVolume, error) {
	persistentVolume := &corev1.PersistentVolume{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.CoreV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), PersistentVolumeKind.Scoped).
		Resource(PersistentVolumeResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(persistentVolume)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   PersistentVolumeKind.Group,
			Version: PersistentVolumeKind.Version,
			Kind:    PersistentVolumeKind.Kind,
		}, persistentVolume.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    PersistentVolumeKind.Group,
				Resource: PersistentVolumeResource.Name,
			}, name)
		}
	}
	return NewPersistentVolume(persistentVolume, c.Client), nil
}

func (c *persistentVolumesReader) List() ([]*PersistentVolume, error) {
	list := &corev1.PersistentVolumeList{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.CoreV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), PersistentVolumeKind.Scoped).
		Resource(PersistentVolumeResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*PersistentVolume, 0, len(list.Items))
	for _, persistentVolume := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   PersistentVolumeKind.Group,
			Version: PersistentVolumeKind.Version,
			Kind:    PersistentVolumeKind.Kind,
		}, persistentVolume.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := persistentVolume
			results = append(results, NewPersistentVolume(&copy, c.Client))
		}
	}
	return results, nil
}

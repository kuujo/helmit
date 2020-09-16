package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

type JobsReader interface {
	Get(name string) (*Job, error)
	List() ([]*Job, error)
}

func NewJobsReader(client resource.Client, filter resource.Filter) JobsReader {
	return &jobsReader{
		Client: client,
		filter: filter,
	}
}

type jobsReader struct {
	resource.Client
	filter resource.Filter
}

func (c *jobsReader) Get(name string) (*Job, error) {
	job := &batchv1.Job{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.BatchV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), JobKind.Scoped).
		Resource(JobResource.Name).
		Name(name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(job)
	if err != nil {
		return nil, err
	} else {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   JobKind.Group,
			Version: JobKind.Version,
			Kind:    JobKind.Kind,
		}, job.ObjectMeta)
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.NewNotFound(schema.GroupResource{
				Group:    JobKind.Group,
				Resource: JobResource.Name,
			}, name)
		}
	}
	return NewJob(job, c.Client), nil
}

func (c *jobsReader) List() ([]*Job, error) {
	list := &batchv1.JobList{}
	client, err := kubernetes.NewForConfig(c.Config())
	if err != nil {
		return nil, err
	}
	err = client.BatchV1().
		RESTClient().
		Get().
		NamespaceIfScoped(c.Namespace(), JobKind.Scoped).
		Resource(JobResource.Name).
		VersionedParams(&metav1.ListOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Into(list)
	if err != nil {
		return nil, err
	}

	results := make([]*Job, 0, len(list.Items))
	for _, job := range list.Items {
		ok, err := c.filter(metav1.GroupVersionKind{
			Group:   JobKind.Group,
			Version: JobKind.Version,
			Kind:    JobKind.Kind,
		}, job.ObjectMeta)
		if err != nil {
			return nil, err
		} else if ok {
			copy := job
			results = append(results, NewJob(&copy, c.Client))
		}
	}
	return results, nil
}

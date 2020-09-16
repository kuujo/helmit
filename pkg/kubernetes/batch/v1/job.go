package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var JobKind = resource.Kind{
	Group:   "batch",
	Version: "v1",
	Kind:    "Job",
	Scoped:  true,
}

var JobResource = resource.Type{
	Kind: JobKind,
	Name: "jobs",
}

func NewJob(job *batchv1.Job, client resource.Client) *Job {
	return &Job{
		Resource: resource.NewResource(job.ObjectMeta, JobKind, client),
		Object:   job,
	}
}

type Job struct {
	*resource.Resource
	Object *batchv1.Job
}

func (r *Job) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.BatchV1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, JobKind.Scoped).
		Resource(JobResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

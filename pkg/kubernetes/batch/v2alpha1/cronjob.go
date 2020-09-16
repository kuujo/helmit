package v2alpha1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var CronJobKind = resource.Kind{
	Group:   "batch",
	Version: "v2alpha1",
	Kind:    "CronJob",
	Scoped:  true,
}

var CronJobResource = resource.Type{
	Kind: CronJobKind,
	Name: "cronjobs",
}

func NewCronJob(cronJob *batchv2alpha1.CronJob, client resource.Client) *CronJob {
	return &CronJob{
		Resource: resource.NewResource(cronJob.ObjectMeta, CronJobKind, client),
		Object:   cronJob,
	}
}

type CronJob struct {
	*resource.Resource
	Object *batchv2alpha1.CronJob
}

func (r *CronJob) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.BatchV2alpha1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, CronJobKind.Scoped).
		Resource(CronJobResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}

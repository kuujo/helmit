package kubernetes

import (
	admissionregistrationv1 "github.com/onosproject/helmit/pkg/kubernetes/admissionregistration/v1"
	apiextensionsv1 "github.com/onosproject/helmit/pkg/kubernetes/apiextensions/v1"
	apiextensionsv1beta1 "github.com/onosproject/helmit/pkg/kubernetes/apiextensions/v1beta1"
	appsv1 "github.com/onosproject/helmit/pkg/kubernetes/apps/v1"
	appsv1beta1 "github.com/onosproject/helmit/pkg/kubernetes/apps/v1beta1"
	batchv1 "github.com/onosproject/helmit/pkg/kubernetes/batch/v1"
	batchv1beta1 "github.com/onosproject/helmit/pkg/kubernetes/batch/v1beta1"
	batchv2alpha1 "github.com/onosproject/helmit/pkg/kubernetes/batch/v2alpha1"
	"github.com/onosproject/helmit/pkg/kubernetes/config"
	corev1 "github.com/onosproject/helmit/pkg/kubernetes/core/v1"
	extensionsv1beta1 "github.com/onosproject/helmit/pkg/kubernetes/extensions/v1beta1"
	networkingv1beta1 "github.com/onosproject/helmit/pkg/kubernetes/networking/v1beta1"
	policyv1beta1 "github.com/onosproject/helmit/pkg/kubernetes/policy/v1beta1"
	rbacv1 "github.com/onosproject/helmit/pkg/kubernetes/rbac/v1"
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	storagev1 "github.com/onosproject/helmit/pkg/kubernetes/storage/v1"
	helmkube "helm.sh/helm/v3/pkg/kube"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// New returns a new Kubernetes client for the current namespace
func New() (Client, error) {
	return NewForNamespace(config.GetNamespaceFromEnv())
}

// NewOrDie returns a new Kubernetes client for the current namespace
func NewOrDie() Client {
	client, err := New()
	if err != nil {
		panic(err)
	}
	return client
}

// NewForNamespace returns a new Kubernetes client for the given namespace
func NewForNamespace(namespace string) (Client, error) {
	kubernetesConfig, err := config.GetRestConfig()
	if err != nil {
		return nil, err
	}
	kubernetesClient, err := kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		return nil, err
	}
	return &client{
		namespace: namespace,
		config:    kubernetesConfig,
		client:    kubernetesClient,
		filter:    resource.NoFilter,
	}, nil
}

// NewForNamespaceOrDie returns a new Kubernetes client for the given namespace
func NewForNamespaceOrDie(namespace string) Client {
	client, err := NewForNamespace(namespace)
	if err != nil {
		panic(err)
	}
	return client
}

// Client is a Kubernetes client
type Client interface {
	// Namespace returns the client namespace
	Namespace() string

	// Config returns the Kubernetes REST client configuration
	Config() *rest.Config

	// Clientset returns the client's Clientset
	Clientset() *kubernetes.Clientset
	AdmissionregistrationV1() admissionregistrationv1.Client
	ApiextensionsV1() apiextensionsv1.Client
	ApiextensionsV1beta1() apiextensionsv1beta1.Client
	AppsV1() appsv1.Client
	AppsV1beta1() appsv1beta1.Client
	BatchV1() batchv1.Client
	BatchV1beta1() batchv1beta1.Client
	BatchV2alpha1() batchv2alpha1.Client
	ExtensionsV1beta1() extensionsv1beta1.Client
	NetworkingV1beta1() networkingv1beta1.Client
	PolicyV1beta1() policyv1beta1.Client
	RbacV1() rbacv1.Client
	StorageV1() storagev1.Client
	CoreV1() corev1.Client
}

// NewForResources returns a new Kubernetes client for the given resources
func NewForResources(namespace string, resources helmkube.ResourceList) (Client, error) {
	kubernetesConfig, err := config.GetRestConfig()
	if err != nil {
		return nil, err
	}
	kubernetesClient, err := kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		return nil, err
	}
	parentClient := &client{
		namespace: namespace,
		config:    kubernetesConfig,
		client:    kubernetesClient,
		filter:    resource.NoFilter,
	}
	return &client{
		namespace: namespace,
		config:    kubernetesConfig,
		client:    kubernetesClient,
		filter: func(kind metav1.GroupVersionKind, meta metav1.ObjectMeta) (bool, error) {
			return filterResources(parentClient, resources, kind, meta)
		},
	}, nil
}

// NewForResourcesOrDie returns a new Kubernetes client for the given release
func NewForResourcesOrDie(namespace string, resources helmkube.ResourceList) Client {
	client, err := NewForResources(namespace, resources)
	if err != nil {
		panic(err)
	}
	return client
}

func filterResources(client resource.Client, resources helmkube.ResourceList, kind metav1.GroupVersionKind, meta metav1.ObjectMeta) (bool, error) {
	for _, resource := range resources {
		resourceKind := resource.Object.GetObjectKind().GroupVersionKind()
		if resourceKind.Group == kind.Group &&
			resourceKind.Version == kind.Version &&
			resourceKind.Kind == kind.Kind &&
			resource.Namespace == meta.Namespace &&
			resource.Name == meta.Name {
			return true, nil
		}
	}
	return filterOwners(client, resources, kind, meta)
}

func filterOwners(client resource.Client, resources helmkube.ResourceList, kind metav1.GroupVersionKind, meta metav1.ObjectMeta) (bool, error) {
	for _, owner := range meta.OwnerReferences {
		ok, err := filterOwner(client, resources, owner)
		if ok {
			return true, nil
		} else if err != nil {
			return false, err
		}
	}
	return filterApp(client, resources, kind, meta)
}

func filterOwner(client resource.Client, resources helmkube.ResourceList, owner metav1.OwnerReference) (bool, error) {
	for _, resource := range resources {
		resourceKind := resource.Object.GetObjectKind().GroupVersionKind()
		if resourceKind.Kind == owner.Kind &&
			resourceKind.GroupVersion().String() == owner.APIVersion &&
			resource.Name == owner.Name {
			return true, nil
		}
	}

	switch owner.APIVersion {
	case "admissionregistration.k8s.io/v1":
		switch owner.Kind {
		case "MutatingWebhookConfiguration":
			mutatingWebhookConfigurationClient := admissionregistrationv1.NewMutatingWebhookConfigurationsReader(client, resource.NoFilter)
			mutatingWebhookConfiguration, err := mutatingWebhookConfigurationClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   admissionregistrationv1.MutatingWebhookConfigurationKind.Group,
					Version: admissionregistrationv1.MutatingWebhookConfigurationKind.Version,
					Kind:    admissionregistrationv1.MutatingWebhookConfigurationKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, mutatingWebhookConfiguration.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "ValidatingWebhookConfiguration":
			validatingWebhookConfigurationClient := admissionregistrationv1.NewValidatingWebhookConfigurationsReader(client, resource.NoFilter)
			validatingWebhookConfiguration, err := validatingWebhookConfigurationClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   admissionregistrationv1.ValidatingWebhookConfigurationKind.Group,
					Version: admissionregistrationv1.ValidatingWebhookConfigurationKind.Version,
					Kind:    admissionregistrationv1.ValidatingWebhookConfigurationKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, validatingWebhookConfiguration.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	case "apiextensions.k8s.io/v1":
		switch owner.Kind {
		case "CustomResourceDefinition":
			customResourceDefinitionClient := apiextensionsv1.NewCustomResourceDefinitionsReader(client, resource.NoFilter)
			customResourceDefinition, err := customResourceDefinitionClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   apiextensionsv1.CustomResourceDefinitionKind.Group,
					Version: apiextensionsv1.CustomResourceDefinitionKind.Version,
					Kind:    apiextensionsv1.CustomResourceDefinitionKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, customResourceDefinition.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	case "apiextensions.k8s.io/v1beta1":
		switch owner.Kind {
		case "CustomResourceDefinition":
			customResourceDefinitionClient := apiextensionsv1beta1.NewCustomResourceDefinitionsReader(client, resource.NoFilter)
			customResourceDefinition, err := customResourceDefinitionClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   apiextensionsv1beta1.CustomResourceDefinitionKind.Group,
					Version: apiextensionsv1beta1.CustomResourceDefinitionKind.Version,
					Kind:    apiextensionsv1beta1.CustomResourceDefinitionKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, customResourceDefinition.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	case "apps/v1":
		switch owner.Kind {
		case "DaemonSet":
			daemonSetClient := appsv1.NewDaemonSetsReader(client, resource.NoFilter)
			daemonSet, err := daemonSetClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1.DaemonSetKind.Group,
					Version: appsv1.DaemonSetKind.Version,
					Kind:    appsv1.DaemonSetKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, daemonSet.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "Deployment":
			deploymentClient := appsv1.NewDeploymentsReader(client, resource.NoFilter)
			deployment, err := deploymentClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1.DeploymentKind.Group,
					Version: appsv1.DeploymentKind.Version,
					Kind:    appsv1.DeploymentKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, deployment.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "ReplicaSet":
			replicaSetClient := appsv1.NewReplicaSetsReader(client, resource.NoFilter)
			replicaSet, err := replicaSetClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1.ReplicaSetKind.Group,
					Version: appsv1.ReplicaSetKind.Version,
					Kind:    appsv1.ReplicaSetKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, replicaSet.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "StatefulSet":
			statefulSetClient := appsv1.NewStatefulSetsReader(client, resource.NoFilter)
			statefulSet, err := statefulSetClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1.StatefulSetKind.Group,
					Version: appsv1.StatefulSetKind.Version,
					Kind:    appsv1.StatefulSetKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, statefulSet.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	case "apps/v1beta1":
		switch owner.Kind {
		case "Deployment":
			deploymentClient := appsv1beta1.NewDeploymentsReader(client, resource.NoFilter)
			deployment, err := deploymentClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1beta1.DeploymentKind.Group,
					Version: appsv1beta1.DeploymentKind.Version,
					Kind:    appsv1beta1.DeploymentKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, deployment.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "StatefulSet":
			statefulSetClient := appsv1beta1.NewStatefulSetsReader(client, resource.NoFilter)
			statefulSet, err := statefulSetClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1beta1.StatefulSetKind.Group,
					Version: appsv1beta1.StatefulSetKind.Version,
					Kind:    appsv1beta1.StatefulSetKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, statefulSet.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	case "batch/v1":
		switch owner.Kind {
		case "Job":
			jobClient := batchv1.NewJobsReader(client, resource.NoFilter)
			job, err := jobClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   batchv1.JobKind.Group,
					Version: batchv1.JobKind.Version,
					Kind:    batchv1.JobKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, job.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	case "batch/v1beta1":
		switch owner.Kind {
		case "CronJob":
			cronJobClient := batchv1beta1.NewCronJobsReader(client, resource.NoFilter)
			cronJob, err := cronJobClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   batchv1beta1.CronJobKind.Group,
					Version: batchv1beta1.CronJobKind.Version,
					Kind:    batchv1beta1.CronJobKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, cronJob.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	case "batch/v2alpha1":
		switch owner.Kind {
		case "CronJob":
			cronJobClient := batchv2alpha1.NewCronJobsReader(client, resource.NoFilter)
			cronJob, err := cronJobClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   batchv2alpha1.CronJobKind.Group,
					Version: batchv2alpha1.CronJobKind.Version,
					Kind:    batchv2alpha1.CronJobKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, cronJob.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	case "extensions/v1beta1":
		switch owner.Kind {
		case "Ingress":
			ingressClient := extensionsv1beta1.NewIngressesReader(client, resource.NoFilter)
			ingress, err := ingressClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   extensionsv1beta1.IngressKind.Group,
					Version: extensionsv1beta1.IngressKind.Version,
					Kind:    extensionsv1beta1.IngressKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, ingress.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	case "networking/v1beta1":
		switch owner.Kind {
		case "Ingress":
			ingressClient := networkingv1beta1.NewIngressesReader(client, resource.NoFilter)
			ingress, err := ingressClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   networkingv1beta1.IngressKind.Group,
					Version: networkingv1beta1.IngressKind.Version,
					Kind:    networkingv1beta1.IngressKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, ingress.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	case "policy/v1beta1":
		switch owner.Kind {
		case "PodDisruptionBudget":
			podDisruptionBudgetClient := policyv1beta1.NewPodDisruptionBudgetsReader(client, resource.NoFilter)
			podDisruptionBudget, err := podDisruptionBudgetClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   policyv1beta1.PodDisruptionBudgetKind.Group,
					Version: policyv1beta1.PodDisruptionBudgetKind.Version,
					Kind:    policyv1beta1.PodDisruptionBudgetKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, podDisruptionBudget.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "PodSecurityPolicy":
			podSecurityPolicyClient := policyv1beta1.NewPodSecurityPoliciesReader(client, resource.NoFilter)
			podSecurityPolicy, err := podSecurityPolicyClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   policyv1beta1.PodSecurityPolicyKind.Group,
					Version: policyv1beta1.PodSecurityPolicyKind.Version,
					Kind:    policyv1beta1.PodSecurityPolicyKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, podSecurityPolicy.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	case "rbac.authorization.k8s.io/v1":
		switch owner.Kind {
		case "ClusterRole":
			clusterRoleClient := rbacv1.NewClusterRolesReader(client, resource.NoFilter)
			clusterRole, err := clusterRoleClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   rbacv1.ClusterRoleKind.Group,
					Version: rbacv1.ClusterRoleKind.Version,
					Kind:    rbacv1.ClusterRoleKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, clusterRole.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "ClusterRoleBinding":
			clusterRoleBindingClient := rbacv1.NewClusterRoleBindingsReader(client, resource.NoFilter)
			clusterRoleBinding, err := clusterRoleBindingClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   rbacv1.ClusterRoleBindingKind.Group,
					Version: rbacv1.ClusterRoleBindingKind.Version,
					Kind:    rbacv1.ClusterRoleBindingKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, clusterRoleBinding.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "Role":
			roleClient := rbacv1.NewRolesReader(client, resource.NoFilter)
			role, err := roleClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   rbacv1.RoleKind.Group,
					Version: rbacv1.RoleKind.Version,
					Kind:    rbacv1.RoleKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, role.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "RoleBinding":
			roleBindingClient := rbacv1.NewRoleBindingsReader(client, resource.NoFilter)
			roleBinding, err := roleBindingClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   rbacv1.RoleBindingKind.Group,
					Version: rbacv1.RoleBindingKind.Version,
					Kind:    rbacv1.RoleBindingKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, roleBinding.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	case "storage.k8s.io/v1":
		switch owner.Kind {
		case "StorageClass":
			storageClassClient := storagev1.NewStorageClassesReader(client, resource.NoFilter)
			storageClass, err := storageClassClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   storagev1.StorageClassKind.Group,
					Version: storagev1.StorageClassKind.Version,
					Kind:    storagev1.StorageClassKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, storageClass.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	case "/v1":
		switch owner.Kind {
		case "ConfigMap":
			configMapClient := corev1.NewConfigMapsReader(client, resource.NoFilter)
			configMap, err := configMapClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   corev1.ConfigMapKind.Group,
					Version: corev1.ConfigMapKind.Version,
					Kind:    corev1.ConfigMapKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, configMap.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "Endpoints":
			endpointsClient := corev1.NewEndpointsReader(client, resource.NoFilter)
			endpoints, err := endpointsClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   corev1.EndpointsKind.Group,
					Version: corev1.EndpointsKind.Version,
					Kind:    corev1.EndpointsKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, endpoints.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "Namespace":
			namespaceClient := corev1.NewNamespacesReader(client, resource.NoFilter)
			namespace, err := namespaceClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   corev1.NamespaceKind.Group,
					Version: corev1.NamespaceKind.Version,
					Kind:    corev1.NamespaceKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, namespace.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "Node":
			nodeClient := corev1.NewNodesReader(client, resource.NoFilter)
			node, err := nodeClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   corev1.NodeKind.Group,
					Version: corev1.NodeKind.Version,
					Kind:    corev1.NodeKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, node.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "PersistentVolume":
			persistentVolumeClient := corev1.NewPersistentVolumesReader(client, resource.NoFilter)
			persistentVolume, err := persistentVolumeClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   corev1.PersistentVolumeKind.Group,
					Version: corev1.PersistentVolumeKind.Version,
					Kind:    corev1.PersistentVolumeKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, persistentVolume.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "PersistentVolumeClaim":
			persistentVolumeClaimClient := corev1.NewPersistentVolumeClaimsReader(client, resource.NoFilter)
			persistentVolumeClaim, err := persistentVolumeClaimClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   corev1.PersistentVolumeClaimKind.Group,
					Version: corev1.PersistentVolumeClaimKind.Version,
					Kind:    corev1.PersistentVolumeClaimKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, persistentVolumeClaim.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "Pod":
			podClient := corev1.NewPodsReader(client, resource.NoFilter)
			pod, err := podClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   corev1.PodKind.Group,
					Version: corev1.PodKind.Version,
					Kind:    corev1.PodKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, pod.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "PodTemplate":
			podTemplateClient := corev1.NewPodTemplatesReader(client, resource.NoFilter)
			podTemplate, err := podTemplateClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   corev1.PodTemplateKind.Group,
					Version: corev1.PodTemplateKind.Version,
					Kind:    corev1.PodTemplateKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, podTemplate.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "Secret":
			secretClient := corev1.NewSecretsReader(client, resource.NoFilter)
			secret, err := secretClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   corev1.SecretKind.Group,
					Version: corev1.SecretKind.Version,
					Kind:    corev1.SecretKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, secret.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		case "Service":
			serviceClient := corev1.NewServicesReader(client, resource.NoFilter)
			service, err := serviceClient.Get(owner.Name)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   corev1.ServiceKind.Group,
					Version: corev1.ServiceKind.Version,
					Kind:    corev1.ServiceKind.Kind,
				}
				ok, err := filterResources(client, resources, groupVersionKind, service.Object.ObjectMeta)
				if ok {
					return true, nil
				} else if err != nil {
					return false, err
				}
			}
		}
	}
	return false, nil
}

func filterApp(client resource.Client, resources helmkube.ResourceList, kind metav1.GroupVersionKind, meta metav1.ObjectMeta) (bool, error) {
	if isSameKind(kind, corev1.PodKind) {
		instance, ok := meta.Labels["app.kubernetes.io/instance"]
		if ok {
			daemonSetClient := appsv1.NewDaemonSetsReader(client, resource.NoFilter)
			daemonSet, err := daemonSetClient.Get(instance)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1.DaemonSetKind.Group,
					Version: appsv1.DaemonSetKind.Version,
					Kind:    appsv1.DaemonSetKind.Kind,
				}
				return filterResources(client, resources, groupVersionKind, daemonSet.Object.ObjectMeta)
			}
		}
	}
	if isSameKind(kind, appsv1.ReplicaSetKind) {
		instance, ok := meta.Labels["app.kubernetes.io/instance"]
		if ok {
			deploymentClient := appsv1.NewDeploymentsReader(client, resource.NoFilter)
			deployment, err := deploymentClient.Get(instance)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1.DeploymentKind.Group,
					Version: appsv1.DeploymentKind.Version,
					Kind:    appsv1.DeploymentKind.Kind,
				}
				return filterResources(client, resources, groupVersionKind, deployment.Object.ObjectMeta)
			}
		}
	}
	if isSameKind(kind, corev1.PodKind) {
		instance, ok := meta.Labels["app.kubernetes.io/instance"]
		if ok {
			replicaSetClient := appsv1.NewReplicaSetsReader(client, resource.NoFilter)
			replicaSet, err := replicaSetClient.Get(instance)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1.ReplicaSetKind.Group,
					Version: appsv1.ReplicaSetKind.Version,
					Kind:    appsv1.ReplicaSetKind.Kind,
				}
				return filterResources(client, resources, groupVersionKind, replicaSet.Object.ObjectMeta)
			}
		}
	}
	if isSameKind(kind, corev1.PodKind) {
		instance, ok := meta.Labels["app.kubernetes.io/instance"]
		if ok {
			statefulSetClient := appsv1.NewStatefulSetsReader(client, resource.NoFilter)
			statefulSet, err := statefulSetClient.Get(instance)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1.StatefulSetKind.Group,
					Version: appsv1.StatefulSetKind.Version,
					Kind:    appsv1.StatefulSetKind.Kind,
				}
				return filterResources(client, resources, groupVersionKind, statefulSet.Object.ObjectMeta)
			}
		}
	}
	if isSameKind(kind, appsv1.ReplicaSetKind) {
		instance, ok := meta.Labels["app.kubernetes.io/instance"]
		if ok {
			deploymentClient := appsv1beta1.NewDeploymentsReader(client, resource.NoFilter)
			deployment, err := deploymentClient.Get(instance)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1beta1.DeploymentKind.Group,
					Version: appsv1beta1.DeploymentKind.Version,
					Kind:    appsv1beta1.DeploymentKind.Kind,
				}
				return filterResources(client, resources, groupVersionKind, deployment.Object.ObjectMeta)
			}
		}
	}
	if isSameKind(kind, corev1.PodKind) {
		instance, ok := meta.Labels["app.kubernetes.io/instance"]
		if ok {
			deploymentClient := appsv1beta1.NewDeploymentsReader(client, resource.NoFilter)
			deployment, err := deploymentClient.Get(instance)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1beta1.DeploymentKind.Group,
					Version: appsv1beta1.DeploymentKind.Version,
					Kind:    appsv1beta1.DeploymentKind.Kind,
				}
				return filterResources(client, resources, groupVersionKind, deployment.Object.ObjectMeta)
			}
		}
	}
	if isSameKind(kind, appsv1.ReplicaSetKind) {
		instance, ok := meta.Labels["app.kubernetes.io/instance"]
		if ok {
			statefulSetClient := appsv1beta1.NewStatefulSetsReader(client, resource.NoFilter)
			statefulSet, err := statefulSetClient.Get(instance)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1beta1.StatefulSetKind.Group,
					Version: appsv1beta1.StatefulSetKind.Version,
					Kind:    appsv1beta1.StatefulSetKind.Kind,
				}
				return filterResources(client, resources, groupVersionKind, statefulSet.Object.ObjectMeta)
			}
		}
	}
	if isSameKind(kind, corev1.PodKind) {
		instance, ok := meta.Labels["app.kubernetes.io/instance"]
		if ok {
			statefulSetClient := appsv1beta1.NewStatefulSetsReader(client, resource.NoFilter)
			statefulSet, err := statefulSetClient.Get(instance)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   appsv1beta1.StatefulSetKind.Group,
					Version: appsv1beta1.StatefulSetKind.Version,
					Kind:    appsv1beta1.StatefulSetKind.Kind,
				}
				return filterResources(client, resources, groupVersionKind, statefulSet.Object.ObjectMeta)
			}
		}
	}
	if isSameKind(kind, corev1.EndpointsKind) {
		instance, ok := meta.Labels["app.kubernetes.io/instance"]
		if ok {
			serviceClient := corev1.NewServicesReader(client, resource.NoFilter)
			service, err := serviceClient.Get(instance)
			if err != nil && !errors.IsNotFound(err) {
				return false, err
			} else if err == nil {
				groupVersionKind := metav1.GroupVersionKind{
					Group:   corev1.ServiceKind.Group,
					Version: corev1.ServiceKind.Version,
					Kind:    corev1.ServiceKind.Kind,
				}
				return filterResources(client, resources, groupVersionKind, service.Object.ObjectMeta)
			}
		}
	}
	return false, nil
}

func isSameKind(groupVersionKind metav1.GroupVersionKind, kind resource.Kind) bool {
	return groupVersionKind.Group == kind.Group &&
		groupVersionKind.Version == kind.Version &&
		groupVersionKind.Kind == kind.Kind
}

type client struct {
	namespace string
	config    *rest.Config
	client    *kubernetes.Clientset
	filter    resource.Filter
}

func (c *client) Namespace() string {
	return c.namespace
}

func (c *client) Config() *rest.Config {
	return c.config
}

func (c *client) Clientset() *kubernetes.Clientset {
	return c.client
}
func (c *client) AdmissionregistrationV1() admissionregistrationv1.Client {
	return admissionregistrationv1.NewClient(c, c.filter)
}

func (c *client) ApiextensionsV1() apiextensionsv1.Client {
	return apiextensionsv1.NewClient(c, c.filter)
}

func (c *client) ApiextensionsV1beta1() apiextensionsv1beta1.Client {
	return apiextensionsv1beta1.NewClient(c, c.filter)
}

func (c *client) AppsV1() appsv1.Client {
	return appsv1.NewClient(c, c.filter)
}

func (c *client) AppsV1beta1() appsv1beta1.Client {
	return appsv1beta1.NewClient(c, c.filter)
}

func (c *client) BatchV1() batchv1.Client {
	return batchv1.NewClient(c, c.filter)
}

func (c *client) BatchV1beta1() batchv1beta1.Client {
	return batchv1beta1.NewClient(c, c.filter)
}

func (c *client) BatchV2alpha1() batchv2alpha1.Client {
	return batchv2alpha1.NewClient(c, c.filter)
}

func (c *client) ExtensionsV1beta1() extensionsv1beta1.Client {
	return extensionsv1beta1.NewClient(c, c.filter)
}

func (c *client) NetworkingV1beta1() networkingv1beta1.Client {
	return networkingv1beta1.NewClient(c, c.filter)
}

func (c *client) PolicyV1beta1() policyv1beta1.Client {
	return policyv1beta1.NewClient(c, c.filter)
}

func (c *client) RbacV1() rbacv1.Client {
	return rbacv1.NewClient(c, c.filter)
}

func (c *client) StorageV1() storagev1.Client {
	return storagev1.NewClient(c, c.filter)
}

func (c *client) CoreV1() corev1.Client {
	return corev1.NewClient(c, c.filter)
}

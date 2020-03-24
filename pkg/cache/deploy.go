package cache

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/go-utils/kubeerrutils"
	"github.com/solo-io/wasme/pkg/version"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// the cache deployer deploys an instance of
// the cache to Kubernetes.
type Deployer interface {
	// ensures the desired cache is deployed.
	// if a cache component already exists, it is updated
	EnsureCache() error
}

var (
	CacheName            = "wasme-cache"
	CacheNamespace       = "wasme"
	CacheImageRepository = "quay.io/solo-io/wasme"
	CacheImageTag        = version.Version
	ImagesKey            = "images"
	DefaultCacheArgs     = func(namespace string) []string {
		return []string{
			"cache",
			"--directory",
			"/var/local/lib/wasme-cache",
			"--ref-file",
			"/etc/wasme-cache/images.txt",
			"--cache-ns",
			namespace,
		}
	}
)

type deployer struct {
	kube       kubernetes.Interface
	namespace  string
	name       string
	image      string
	pullPolicy v1.PullPolicy
	args       []string
	logger     *logrus.Entry
}

func NewDeployer(kube kubernetes.Interface, namespace, name string, imageRepo, imageTag string, args []string, pullPolicy v1.PullPolicy) *deployer {
	if namespace == "" {
		namespace = CacheNamespace
	}
	if name == "" {
		name = CacheName
	}
	if imageRepo == "" {
		imageRepo = CacheImageRepository
	}
	if imageTag == "" {
		imageTag = CacheImageTag
	}
	if args == nil {
		args = DefaultCacheArgs(namespace)
	}
	image := imageRepo + ":" + imageTag
	return &deployer{
		kube:       kube,
		namespace:  namespace,
		name:       name,
		image:      image,
		args:       args,
		pullPolicy: pullPolicy,
		logger: logrus.WithFields(logrus.Fields{
			"cache": name + "." + namespace,
			"image": image,
		},
		)}
}

func (d *deployer) EnsureCache() error {
	if err := d.createNamespaceIfNotExist(); err != nil {
		return errors.Wrap(err, "ensuring namespace")
	}

	if err := d.createServiceAccountIfNotExist(); err != nil {
		return errors.Wrap(err, "ensuring service acct")
	}

	role, roleBinding := MakeRbac(d.name, d.namespace)

	if err := d.createOrUpdateCacheRole(role); err != nil {
		return errors.Wrap(err, "ensuring role")
	}

	if err := d.createOrUpdateCacheRolebinding(roleBinding); err != nil {
		return errors.Wrap(err, "ensuring rolebinding")
	}

	if err := d.createOrUpdateDeployment(); err != nil {
		return errors.Wrap(err, "ensuring deployment")
	}
	return nil
}

func (d *deployer) createNamespaceIfNotExist() error {
	_, err := d.kube.CoreV1().Namespaces().Create(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.namespace,
		},
	})
	// ignore already exists err
	if err != nil {
		if kubeerrutils.IsAlreadyExists(err) {
			d.logger.Info("cache namespace already exists")
			return nil
		}
		return err
	}
	d.logger.Info("cache namespace created")
	return nil
}

func (d *deployer) createServiceAccountIfNotExist() error {
	svcAcct := MakeServiceAccount(d.name, d.namespace)
	_, err := d.kube.CoreV1().ServiceAccounts(d.namespace).Create(svcAcct)
	// ignore already exists err
	if err != nil {
		if kubeerrutils.IsAlreadyExists(err) {
			d.logger.Info("cache service account already exists")
			return nil
		}
		return err
	}
	d.logger.Info("cache service account created")
	return err
}

func (d *deployer) createOrUpdateCacheRole(role *rbacv1.Role) error {

	_, err := d.kube.RbacV1().Roles(d.namespace).Create(role)
	// update on already exists err
	if err != nil {
		if !kubeerrutils.IsAlreadyExists(err) {
			return err
		}
		existing, err := d.kube.RbacV1().Roles(d.namespace).Get(role.Name, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to get existing cache role")
		}

		existing.Rules = role.Rules

		_, err = d.kube.RbacV1().Roles(d.namespace).Update(existing)
		if err != nil {
			return err
		}

		d.logger.Info("cache role updated")

		return nil
	}

	d.logger.Info("cache role created")

	return nil
}

func (d *deployer) createOrUpdateCacheRolebinding(roleBinding *rbacv1.RoleBinding) error {
	_, err := d.kube.RbacV1().RoleBindings(d.namespace).Create(roleBinding)
	// update on already exists err
	if err != nil {
		if !kubeerrutils.IsAlreadyExists(err) {
			return err
		}
		existing, err := d.kube.RbacV1().RoleBindings(d.namespace).Get(roleBinding.Name, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to get existing cache rolebinding")
		}

		existing.Subjects = roleBinding.Subjects
		existing.RoleRef = roleBinding.RoleRef

		_, err = d.kube.RbacV1().RoleBindings(d.namespace).Update(existing)
		if err != nil {
			return err
		}

		d.logger.Info("cache rolebinding updated")

		return nil
	}

	d.logger.Info("cache rolebinding created")

	return nil
}

func (d *deployer) createOrUpdateDeployment() error {
	labels := map[string]string{
		"app": d.name,
	}

	desiredDeployment := MakeDeployment(d.name, d.namespace, d.image, labels, d.args, d.pullPolicy)

	_, err := d.kube.AppsV1().Deployments(d.namespace).Create(desiredDeployment)
	// update on already exists err
	if err != nil {
		if !kubeerrutils.IsAlreadyExists(err) {
			return err
		}
		existing, err := d.kube.AppsV1().Deployments(d.namespace).Get(desiredDeployment.Name, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to get existing cache deployment")
		}

		// TODO: how will this handle immutable fields?
		existing.Spec = desiredDeployment.Spec

		_, err = d.kube.AppsV1().Deployments(d.namespace).Update(existing)
		if err != nil {
			return err
		}

		d.logger.Info("cache deployment updated")

		return nil
	}

	d.logger.Info("cache deployment created")

	return nil
}

func MakeServiceAccount(name, namespace string) *v1.ServiceAccount {
	return &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func MakeRbac(name, namespace string) (*rbacv1.Role, *rbacv1.RoleBinding) {
	meta := metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}
	role := &rbacv1.Role{
		ObjectMeta: meta,
		// creates events
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:     []string{"create"},
				APIGroups: []string{""},
				Resources: []string{"events"},
			},
		},
	}
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: meta,
		Subjects: []rbacv1.Subject{{
			Kind: "ServiceAccount",
			Name: name,
		}},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     name,
		},
	}

	return role, roleBinding
}

func MakeDeployment(name, namespace, image string, labels map[string]string, args []string, pullPolicy v1.PullPolicy) *appsv1.Deployment {
	hostPathType := v1.HostPathDirectoryOrCreate
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					ServiceAccountName: name,
					Volumes:            []v1.Volume{},
					Containers: []v1.Container{{
						Name:            name,
						Image:           image,
						ImagePullPolicy: pullPolicy,
						Args:            args,
						VolumeMounts:    []v1.VolumeMount{},
						Resources: v1.ResourceRequirements{
							Limits: v1.ResourceList{
								v1.ResourceMemory: resource.MustParse("256Mi"),
								v1.ResourceCPU:    resource.MustParse("500m"),
							},
							Requests: v1.ResourceList{
								v1.ResourceMemory: resource.MustParse("128Mi"),
								v1.ResourceCPU:    resource.MustParse("50m"),
							},
						},
					}},
				},
			},
		},
	}
}

// get the cache events for an image.
// used by tests and the istio deployer, not by this package
func GetImageEvents(kube kubernetes.Interface, eventNamespace, image string) ([]v1.Event, error) {
	imageEvents, err := kube.CoreV1().Events(eventNamespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(EventLabels(image)).String(),
	})
	if err != nil {
		return nil, err
	}
	return imageEvents.Items, nil
}

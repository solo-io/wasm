package cache

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/go-utils/kubeerrutils"
	"github.com/solo-io/wasme/pkg/version"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	DefaultCacheArgs     = []string{
		"cache",
		"--directory",
		"/var/local/lib/wasme-cache",
		"--ref-file",
		"/etc/wasme-cache/images.txt",
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
		args = DefaultCacheArgs
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
	if err := d.createConfigMapIfNotExist(); err != nil {
		return errors.Wrap(err, "ensuring configmap")
	}
	if err := d.createOrUpdateDaemonSet(); err != nil {
		return errors.Wrap(err, "ensuring daemonset")
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

func (d *deployer) createConfigMapIfNotExist() error {
	_, err := d.kube.CoreV1().ConfigMaps(d.namespace).Create(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.name,
			Namespace: d.namespace,
		},
		Data: map[string]string{
			ImagesKey: "",
		},
	})
	// ignore already exists err
	if err != nil {
		if kubeerrutils.IsAlreadyExists(err) {
			d.logger.Info("cache configmap already exists")
			return nil
		}
		return err
	}
	d.logger.Info("cache configmap created")
	return err
}

func (d *deployer) createOrUpdateDaemonSet() error {
	labels := map[string]string{
		"app": d.name,
	}

	desiredDaemonSet := MakeDaemonSet(d.name, d.namespace, d.image, labels, d.args, d.pullPolicy)

	_, err := d.kube.AppsV1().DaemonSets(d.namespace).Create(desiredDaemonSet)
	// update on already exists err
	if err != nil {
		if !kubeerrutils.IsAlreadyExists(err) {
			return err
		}
		existing, err := d.kube.AppsV1().DaemonSets(d.namespace).Get(desiredDaemonSet.Name, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to get existing cache daemonset")
		}

		// TODO: how will this handle immutable fields?
		existing.Spec = desiredDaemonSet.Spec

		_, err = d.kube.AppsV1().DaemonSets(d.namespace).Update(existing)
		if err != nil {
			return err
		}

		d.logger.Info("cache daemonset updated")

		return nil
	}

	d.logger.Info("cache daemonset created")

	return nil
}

func MakeDaemonSet(name, namespace, image string, labels map[string]string, args []string, pullPolicy v1.PullPolicy) *appsv1.DaemonSet {
	hostPathType := v1.HostPathDirectoryOrCreate
	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "cache-dir",
							VolumeSource: v1.VolumeSource{
								HostPath: &v1.HostPathVolumeSource{
									Path: "/var/local/lib/wasme-cache",
									Type: &hostPathType,
								},
							},
						},
						{
							Name: "config",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: name,
									},
									Items: []v1.KeyToPath{
										{
											Key:  "images",
											Path: "images.txt",
										},
									},
								},
							},
						},
					},
					Containers: []v1.Container{{
						Name:            name,
						Image:           image,
						ImagePullPolicy: pullPolicy,
						Args:            args,
						VolumeMounts: []v1.VolumeMount{
							{
								MountPath: "/var/local/lib/wasme-cache",
								Name:      "cache-dir",
							},
							{
								MountPath: "/etc/wasme-cache",
								Name:      "config",
							},
						},
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

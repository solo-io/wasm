package main

import (
	"github.com/solo-io/wasme/pkg/cache"
	"github.com/solo-io/wasme/pkg/version"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"

	"github.com/solo-io/autopilot/codegen"
	"github.com/solo-io/autopilot/codegen/model"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	hostPathType := v1.HostPathDirectoryOrCreate

	cmd := &codegen.Command{
		AppName: "wasme",
		Groups: []model.Group{
			{
				ProtoDir: "operator/api",
				GroupVersion: schema.GroupVersion{
					Group:   "wasme.io",
					Version: "v1",
				},
				Module: "github.com/solo-io/wasme",
				Resources: []model.Resource{
					{
						Kind: "FilterDeployment",
						Spec: model.Field{
							Type: "FilterDeploymentSpec",
						},
						Status: &model.Field{
							Type: "FilterDeploymentStatus",
						},
					},
				},
				RenderProtos:     true,
				RenderManifests:  true,
				RenderTypes:      true,
				RenderClients:    false,
				RenderController: true,
				ApiRoot:          "pkg/operator/api",
			},
		},

		Chart: &model.Chart{
			Operators: []model.Operator{
				{
					Name: "wasme-operator",
					Deployment: model.Deployment{
						Image: model.Image{
							Registry:   "quay.io/solo-io",
							Repository: "wasme",
							Tag:        version.Version,
							PullPolicy: v1.PullAlways,
							Build: &model.BuildOptions{
								MainFile: "cmd/main.go",
								Push:     true,
							},
						},
					},
					Rbac: []rbacv1.PolicyRule{
						// api resource
						{
							Verbs:     []string{"get", "list", "watch"},
							APIGroups: []string{"wasme.io"},
							Resources: []string{"filterdeployments"},
						},
						{
							Verbs:     []string{"get", "update"},
							APIGroups: []string{"wasme.io"},
							Resources: []string{"filterdeployments/status"},
						},

						// dependency
						{
							Verbs:     []string{"get", "list", "watch"},
							APIGroups: []string{""},
							Resources: []string{"secrets"},
						},

						// managed resources
						{
							Verbs:     []string{"get", "list", "watch", "update"},
							APIGroups: []string{"apps"},
							Resources: []string{"deployments", "daemonsets"},
						},
						{
							Verbs:     []string{"*"},
							APIGroups: []string{"networking.istio.io"},
							Resources: []string{"envoyfilters"},
						},
						{
							Verbs:     []string{"*"},
							APIGroups: []string{""},
							Resources: []string{"configmaps"},
						},
					},
					Args: []string{
						"operator",
						"--log-level=debug",
					},
				},

				// cache
				{
					Name: "wasme-cache",
					Deployment: model.Deployment{
						Image: model.Image{
							Registry:   "quay.io/solo-io",
							Repository: "wasme",
							Tag:        version.Version,
							PullPolicy: v1.PullAlways,
						},
						UseDaemonSet: true,
					},
					Args: []string{
						"cache",
						"--directory",
						"/var/local/lib/wasme-cache",
						"--ref-file",
						"/etc/wasme-cache/images.txt",
					},
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
										Name: cache.CacheName,
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
					ConfigMaps: []v1.ConfigMap{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: cache.CacheName,
							},
							Data: map[string]string{
								cache.ImagesKey: "",
							},
						},
					},
				},
			},
			Values: nil,
			Data: model.Data{
				ApiVersion:  "v1",
				Description: "",
				Name:        "Wasme Operator",
				Version:     "v0.0.1",
				Home:        "https://docs.solo.io/web-assembly-hub/latest",
				Icon:        "https://raw.githubusercontent.com/solo-io/wasme/master/docs/content/img/logo.png",
				Sources: []string{
					"https://github.com/solo-io/wasme",
				},
			},
		},

		ManifestRoot: "operator/install/kube",
		BuildRoot:    "operator/build",
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}

	log.Printf("executed generation with opts: %v", cmd)
}

package main

import (
	"log"
	"os"

	"github.com/solo-io/wasme/pkg/cache"
	"github.com/solo-io/wasme/pkg/version"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/autopilot/codegen"
	"github.com/solo-io/autopilot/codegen/model"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	pushImage := os.Getenv("IMAGE_PUSH") == "1"

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
				RenderClients:    true,
				RenderController: true,
				ApiRoot:          "pkg/operator/api",
			},
		},

		Chart: &model.Chart{
			Operators: []model.Operator{
				makeOperator(),
				makeCache(),
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

		ManifestRoot: "operator/install",

		Builds: []model.Build{
			{
				MainFile: "cmd/main.go",
				Push:     pushImage,
				Image:    makeImage(),
			},
		},
		BuildRoot: "operator/build",
	}
	log.Printf("generating operator with opts: %v", cmd)

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}

	log.Printf("operator generation successful")
}

var (
	defaultRegistry = "quay.io/solo-io"
)

// cache and operator share same image
func makeImage() model.Image {
	registry := os.Getenv("IMAGE_REGISTRY")
	if registry == "" {
		registry = defaultRegistry
	}
	return model.Image{
		Registry:   registry,
		Repository: "wasme",
		Tag:        version.Version,
		PullPolicy: v1.PullAlways,
	}
}

func makeOperator() model.Operator {
	return model.Operator{
		Name: "wasme-operator",
		Deployment: model.Deployment{
			Image: makeImage(),
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
	}
}

func makeCache() model.Operator {
	// need to take a pointer
	hostPathType := v1.HostPathDirectoryOrCreate
	hostPathTypePtr := &hostPathType

	return model.Operator{
		Name: "wasme-cache",
		Deployment: model.Deployment{
			Image: makeImage(),
			Resources: &v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse("125m"),
					v1.ResourceMemory: resource.MustParse("256Mi"),
				},
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
						Type: hostPathTypePtr,
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
	}
}

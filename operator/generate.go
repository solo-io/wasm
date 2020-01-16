package main

import (
	"github.com/solo-io/wasme/pkg/version"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"log"

	"github.com/solo-io/autopilot/codegen"
	"github.com/solo-io/autopilot/codegen/model"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
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
					Name: "wasme",
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
							Verbs:           []string{"get", "list", "watch"},
							APIGroups:       []string{"wasme.io"},
							Resources:       []string{"filterdeployments"},
						},
						{
							Verbs:           []string{"get", "update"},
							APIGroups:       []string{"wasme.io"},
							Resources:       []string{"filterdeployments/status"},
						},

						// dependency
						{
							Verbs:           []string{"get", "list", "watch"},
							APIGroups:       []string{"apps"},
							Resources:       []string{"deployments", "daemonsets"},
						},

						// managed resource
						{
							Verbs:           []string{"*"},
							APIGroups:       []string{"networking.istio.io"},
							Resources:       []string{"envoyfilters"},
						},
					},
					Args: []string{"operator"},
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

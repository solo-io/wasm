package main

import (
	"github.com/solo-io/wasme/pkg/version"
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
				RenderClients:    true,
				RenderController: true,
				ApiRoot:          "pkg/operator/api",
			},
		},
		ManifestRoot: "operator/install/kube",

		Chart: &model.Chart{
			Operators: []model.Operator{
				{
					Name: "wasme",
					Deployment: model.Deployment{
						Image: model.Image{
							Tag:        version.Version,
							Repository: "wasme",
							Registry:   "quay.io/solo-io",
							PullPolicy: "IfNotPresent",
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
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}

	log.Printf("executed generation with opts: %v", cmd)
}

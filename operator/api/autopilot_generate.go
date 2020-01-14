package main

import (
	"log"

	"github.com/solo-io/autopilot/codegen"
	"github.com/solo-io/autopilot/codegen/model"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	cmd := &codegen.Command{
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
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}

	log.Printf("executed generation with opts: %v", cmd)
}

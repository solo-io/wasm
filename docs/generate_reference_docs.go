package main

//go:generate go run generate_reference_docs.go

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogo/protobuf/proto"
	plugin_gogo "github.com/gogo/protobuf/protoc-gen-gogo/plugin"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
	gendoc "github.com/pseudomuto/protoc-gen-doc"
	"github.com/pseudomuto/protoc-gen-doc/parser"
	"github.com/solo-io/autopilot/codegen/util"
	"github.com/solo-io/solo-kit/pkg/code-generator/collector"

	"github.com/solo-io/wasme/pkg/cmd"

	"github.com/solo-io/go-utils/clidoc"
)

var (
	moduleRoot       = util.GetModuleRoot()
	referenceDir     = filepath.Join(moduleRoot, "docs", "content", "reference")
	protoDocTemplate = filepath.Join(moduleRoot, "docs", "proto_docs_template.tmpl")

	cliDocsDir = filepath.Join(referenceDir, "cli")
	cliIndex   = `
---
title: "Command-Line Reference"
description: | 
  Detailed descriptions and options for working with the wasme CLI, which is used to build, manage, and deploy WASM Filters. 
weight: 2
---

This section contains generated reference documentation for the ` + "`" + `wasme` + "`" + ` CLI.

{{% children description="true" %}}

`

	operatorApiDir       = filepath.Join(moduleRoot, "operator", "api", "wasme", "v1")
	operatorDocsTemplate = protoDocTemplate
	operatorDocsDir      = filepath.Join(referenceDir, "operator")
	operatorIndex        = `
---
title: "Operator API Reference"
description: | 
  This section contains the API Specification for the CRDs used by the wasme Kubernetes Operator. The operator is used to deploy Envoy filters to Istio using declarative configuration.
weight: 4
---

These docs describe the ` + "`" + `spec` + "`" + ` and ` + "`" + `status` + "`" + ` of the Wasme Operator's CRD, the FilterDeployment.

{{% children description="true" %}}

`

	imageConfigApiDir       = filepath.Join(moduleRoot, "pkg", "config")
	imageConfigDocsTemplate = protoDocTemplate
	imageConfigDocsDir      = filepath.Join(referenceDir, "image_config")
	imageConfigIndex        = `
---
title: "Wasme Image Config Reference"
description: | 
  For those who build WASM images with wasme, an image config manifest is required by the wasme build commmand. This section contains the API Specification for the Image Config file.
weight: 6
---

This section contains generated reference documentation for the ` + "`" + `wasme` + "`" + ` Filter Image Config.

A config file is packaged at build time with all Envoy WASM Filter Images with metadata and runtime configuration for Wasme.

These docs describe the contents of the Image Config. See https://github.com/solo-io/wasme/blob/master/example/cpp/runtime-config.json
for an example runtime configuration.

{{% children description="true" %}}

`
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := generateCliReference(); err != nil {
		return err
	}
	if err := generateOperatorReference(); err != nil {
		return err
	}
	if err := generateImageConfigReference(); err != nil {
		return err
	}
	return nil
}

func generateCliReference() error {
	// flush directory for idempotence
	os.RemoveAll(cliDocsDir)
	os.MkdirAll(cliDocsDir, 0755)
	app := cmd.Cmd()
	err := clidoc.GenerateCliDocsWithConfig(app, clidoc.Config{
		OutputDir: cliDocsDir,
	})
	if err != nil {
		return errors.Errorf("error generating docs: %s", err)
	}
	return ioutil.WriteFile(filepath.Join(cliDocsDir, "_index.md"), []byte(cliIndex), 0644)
}

func generateOperatorReference() error {
	// flush directory for idempotence
	os.RemoveAll(operatorDocsDir)
	os.MkdirAll(operatorDocsDir, 0755)

	return generateProtoDocs(operatorApiDir, operatorDocsTemplate, operatorDocsDir, operatorIndex)
}

func generateImageConfigReference() error {
	// flush directory for idempotence
	os.RemoveAll(imageConfigDocsDir)
	os.MkdirAll(imageConfigDocsDir, 0755)

	return generateProtoDocs(imageConfigApiDir, imageConfigDocsTemplate, imageConfigDocsDir, imageConfigIndex)
}

func generateProtoDocs(protoDir, templateFile, destDir, indexContents string) error {
	tmpDir, err := ioutil.TempDir("", "proto-docs")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmpDir)

	docsTemplate, err := collectDescriptors(protoDir, tmpDir, nil)
	if err != nil {
		return err
	}

	templateContents, err := ioutil.ReadFile(templateFile)

	tmpl, err := template.New(templateFile).Parse(string(templateContents))
	if err != nil {
		return err
	}

	for _, file := range docsTemplate.Files {
		filename := filepath.Join(destDir, filepath.Base(file.Name))
		filename = strings.TrimSuffix(filename, ".proto") + ".md"
		destFile, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer destFile.Close()
		if err := tmpl.Execute(destFile, file); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(filepath.Join(destDir, "_index.md"), []byte(indexContents), 0644)
}

func collectDescriptors(protoDir, outDir string, customImports []string) (*gendoc.Template, error) {
	descriptors, err := collector.NewCollector(
		customImports,
		[]string{protoDir}, // import the inputs dir
		nil,
		[]string{},
		outDir,
		func(file string) bool {
			return true
		}).CollectDescriptorsFromRoot(protoDir, nil)
	if err != nil {
		return nil, err
	}

	req := &plugin_gogo.CodeGeneratorRequest{}
	for _, file := range descriptors {
		var added bool
		for _, addedFile := range req.GetFileToGenerate() {
			if addedFile == file.GetName() {
				added = true
			}
		}
		if added {
			continue
		}
		req.FileToGenerate = append(req.FileToGenerate, file.GetName())
		req.ProtoFile = append(req.ProtoFile, file.FileDescriptorProto)
	}

	// we have to convert the codegen request from a gogo proto to a golang proto
	// because of incompatibility between the solo kit Collector and the
	// psuedomoto/protoc-doc-gen library:
	golangRequest, err := func() (*plugin_go.CodeGeneratorRequest, error) {
		b, err := proto.Marshal(req)
		if err != nil {
			return nil, err
		}
		var golangReq plugin_go.CodeGeneratorRequest
		if err := proto.Unmarshal(b, &golangReq); err != nil {
			return nil, err
		}
		return &golangReq, nil
	}()

	return gendoc.NewTemplate(parser.ParseCodeRequest(golangRequest)), nil
}

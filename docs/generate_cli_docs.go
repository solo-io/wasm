package main

//go:generate go run generate_cli_docs.go

import (
	"github.com/pkg/errors"
	"github.com/solo-io/autopilot/codegen/util"
	"github.com/solo-io/solo-kit/pkg/code-generator/collector"
	"github.com/solo-io/solo-kit/pkg/code-generator/model"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/solo-io/wasme/pkg/cmd"

	"github.com/solo-io/go-utils/clidoc"
)

var (
	moduleRoot   = util.GetModuleRoot()
	referenceDir = filepath.Join(moduleRoot, "content", "reference")

	cliDocsDir = filepath.Join(referenceDir, "cli")
	cliIndex        = `
---
title: "Command-Line Reference"
weight: 2
---

This section contains generated reference documentation for the `+"`"+`wasme`+"`"+` CLI.

{{% children description="true" %}}

`

	operatorApiDir       = filepath.Join(moduleRoot, "operator", "api", "wasme", "v1")
	operatorDocsTemplate = filepath.Join(operatorApiDir, "proto_docs_template.tmpl")
	operatorDocsDir      = filepath.Join(referenceDir, "operator")
	operatorIndex        = `
---
title: "Operator API Reference"
weight: 4
---

This section contains generated reference documentation for the ` + "`" + `wasme` + "`" + ` Kubernetes Operator.

These docs describe the ` + "`" + `spec` + "`" + ` and ` + "`" + `status` + "`" + ` of Wasme's CRDs.

{{% children description="true" %}}

`

	imageConfigApiDir       = filepath.Join(moduleRoot, "image_config", "pkg", "config")
	imageConfigDocsTemplate = filepath.Join(imageConfigApiDir, "proto_docs_template.tmpl")
	imageConfigDocsDir      = filepath.Join(referenceDir, "image_config")
	imageConfigIndex        = `
---
title: "Wasme Image Config Reference"
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
	generateCliReference()
	generateOperatorReference()
	generateImageConfigReference()
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
	return ioutil.WriteFile(filepath.Join(cliDocsDir, "_index.md"), []byte(operatorIndex), 0644)
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

	descriptors, err := collectDescriptors(protoDir, tmpDir, nil)
	if err != nil {
		return err
	}

	templateContents, err := ioutil.ReadFile(templateFile)

	tmpl, err := template.New(templateFile).Parse(string(templateContents))
	if err != nil {
		return err
	}

	for _, descriptor := range descriptors {
		filename := filepath.Join(destDir, filepath.Base(descriptor.ProtoFilePath))
		destFile, err := os.Create(filename)
		if err != nil {
			return err
		}

		if err := tmpl.Execute(destFile, descriptor); err != nil {
			return err
		}
		if err := destFile.Close(); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(filepath.Join(destDir, "_index.md"), []byte(indexContents), 0644)
}

func collectDescriptors(protoDir, outDir string, customImports []string) ([]*model.DescriptorWithPath, error) {
	return collector.NewCollector(
		customImports,
		[]string{protoDir}, // import the inputs dir
		nil,
		[]string{},
		outDir,
		func(file string) bool {
			return true
		}).CollectDescriptorsFromRoot(protoDir, nil)
}

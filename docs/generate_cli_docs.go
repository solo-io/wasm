package main

//go:generate go run generate_cli_docs.go

import (
	"github.com/solo-io/wasme/pkg/cmd"
	"io/ioutil"
	"log"
	"os"

	"github.com/solo-io/go-utils/clidoc"
)

func main() {
	out := "./content/reference/cli"
	os.RemoveAll(out)
	os.MkdirAll(out, 0755)
	app := cmd.Cmd()
	err := clidoc.GenerateCliDocsWithConfig(app, clidoc.Config{
		OutputDir: out,
	})
	if err != nil {
		log.Fatalf("error generating docs: %s", err)
	}
	err = ioutil.WriteFile(out+"/_index.md", []byte(`
---
title: "Command-Line Reference"
weight: 2
---

This section contains generated reference documentation for the `+"`"+`wasme`+"`"+` CLI.

{{% children description="true" %}}

`), 0644)
	if err != nil {
		log.Fatalf("error writing _index: %s", err)
	}
}

// This file is used to generate embeddable binary data (in Go) from the content of the example directory
// Outputs to pkg/cmd/initialize/cpp_archive_2gobytes.go
package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasme/pkg/util"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

//go:generate go run generate.go

var log = logrus.StandardLogger()

// get the examples dir
var examplesDir = func() string {
	d, err := getCallerDirectory()
	if err != nil {
		log.Fatalf("internal error: failed to get caller directory")
	}
	return d
}()

func getCallerDirectory(skip ...int) (string, error) {
	s := 1
	if len(skip) > 0 {
		s = skip[0] + 1
	}
	_, callerFile, _, ok := runtime.Caller(s)
	if !ok {
		return "", fmt.Errorf("failed to get runtime.Caller")
	}
	callerDir := filepath.Dir(callerFile)

	return filepath.Abs(callerDir)
}

// add to this set to add more example types
// key is the prefix of the variable name
// value is the directory name
var examples = map[string]string{
	"cpp": "cpp",
	"cppIstio1_4": "cpp-istio-1.4",
}

func generateEmbeddedArchiveFile(prefix, dir string) error {
	// tar dir
	archive, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	if err := util.Tar(dir, archive); err != nil {
		return err
	}

	// generate embedded file with 2goarray

	destFile := examplesDir + "/../pkg/cmd/initialize/" + dir + "_archive_2gobytes.go"

	script := fmt.Sprintf("2goarray %sTarBytes initialize < %s | sed 's@// date.*@@g' > %s ", prefix, archive.Name(), destFile)

	genCmd := exec.Command("sh", "-c", script)

	genCmd.Stderr = os.Stderr

	if err := genCmd.Run(); err != nil {
		return err
	}

	return nil
}

func run() error {
	for prefix, dir := range examples {
		if err := generateEmbeddedArchiveFile(prefix, dir); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

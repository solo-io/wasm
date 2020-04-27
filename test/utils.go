package test

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/onsi/ginkgo"

	"github.com/solo-io/autopilot/cli/pkg/utils"
	"github.com/solo-io/autopilot/codegen/util"
	"github.com/solo-io/wasme/pkg/cmd"
)

// split the args from a single line by whitespace
func WasmeCliSplit(argLine string) error {
	log.Printf("arg) wasme: %v", argLine)
	args := strings.Split(argLine, " ")
	return WasmeCli(args...)
}

func WasmeCli(args ...string) error {
	c := cmd.Cmd()
	c.SetArgs(args)
	c.InOrStdin()
	return c.Execute()
}

func RunMake(target string, opts ...func(*exec.Cmd)) error {
	cmd := exec.Command("make", "-B", "-C", filepath.Dir(util.GoModPath()), target)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd.Run()
}

func ApplyFile(file, ns string) error {
	return WithManifest(file, ns, utils.KubectlApply)
}

func DeleteFile(file, ns string) error {
	return WithManifest(file, ns, utils.KubectlDelete)
}

// execute a callback for a manifest relative to the root of the project
func WithManifest(file, ns string, do func(manifest []byte, extraArgs ...string) error) error {
	path := filepath.Join(filepath.Dir(util.GoModPath()), file)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	extraArgs := []string{}
	if ns != "" {
		extraArgs = []string{"-n", ns}
	}
	return do(b, extraArgs...)
}

func GetImageTag() string {
	return GetEnv("FILTER_IMAGE_TAG")
}
func GetBuildImageTag() string {
	return GetEnv("FILTER_BUILD_IMAGE_TAG")
}

func GetEnv(env string) string {
	val := strings.TrimSpace(os.Getenv(env))
	if val == "" {
		ginkgo.Skip("Skipping build/push test. To enable, set " + env + " to the tag to use for the built/pushed image")
	}
	return val
}

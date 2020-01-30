package test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/solo-io/autopilot/cli/pkg/utils"
	"github.com/solo-io/autopilot/codegen/util"
	"github.com/solo-io/wasme/pkg/cmd"
)

// split the args from a single line by whitespace
func WasmeCliSplit(argLine string) error {
	args := strings.Split(argLine, " ")
	return WasmeCli(args...)
}

func WasmeCli(args ...string) error {
	c := cmd.Cmd()
	c.SetArgs(args)
	c.InOrStdin()
	return c.Execute()
}

func RunMake(target string) error {
	cmd := exec.Command("make", "-C", filepath.Dir(util.GoModPath()), target)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
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

package build

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var log = logrus.StandardLogger()

type buildOptions struct {
	sourceDir    string
	outFile      string
	builderImage string
}

func BuildCmd() *cobra.Command {
	var opts buildOptions
	cmd := &cobra.Command{
		Use:   "build SOURCE_DIRECTORY [-o OUTPUT_FILE]",
		Short: "Compile the filter to wasm using Bazel-in-Docker",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("invalid number of arguments")
			}
			opts.sourceDir = args[0]
			return runBuild(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.outFile, "output", "o", "_output/filter.wasm", "path to the output .wasm file. Nonexistent directories will be created.")
	cmd.Flags().StringVarP(&opts.builderImage, "image", "i", "quay.io/solo-io/ee-builder:v1", "Name of the docker image containing the Bazel run instructions. Only be modified if you are an experiejnced user")
	return cmd
}

func runBuild(opts buildOptions) error {
	sourceDir, err := filepath.Abs(opts.sourceDir)
	if err != nil {
		return err
	}
	outFile, err := filepath.Abs(opts.outFile)
	if err != nil {
		return err
	}

	tmpDir, err := ioutil.TempDir("/tmp", "extend-envoy")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	// container paths are currently hard-coded in builder image
	args := []string{
		"run",
		"-v", sourceDir + ":/src/workspace",
		"-v", tmpDir + ":/tmp/build_output",
		"-w", "/src/workspace",
		opts.builderImage,
	}

	log.WithFields(logrus.Fields{
		"args": args,
	}).Info("running bazel-in-docker build...")

	if err := docker(os.Stdout, os.Stderr, args...); err != nil {
		return err
	}

	// filter.wasm currently hard-coded in bazel BUILD file
	tmpFile := filepath.Join(tmpDir, "filter.wasm")

	log.WithFields(logrus.Fields{
		"tmp_file":    tmpFile,
		"output_file": outFile,
	}).Info("moving output file...")

	if err := os.MkdirAll(filepath.Dir(outFile), 0755); err != nil {
		return err
	}

	if err := os.Rename(tmpFile, outFile); err != nil {
		return err
	}

	if err := os.Chmod(outFile, 0644); err != nil{
		return err
	}

	log.WithFields(logrus.Fields{
		"output_file": outFile,
	}).Info("compilation complete!")

	return nil
}

func docker(stdout, stderr io.Writer, args ...string) error {
	return execCmd(stdout, stderr, "docker", args...)
}

func execCmd(stdout, stderr io.Writer, cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	command.Stderr = stderr
	command.Stdout = stdout
	return command.Run()
}

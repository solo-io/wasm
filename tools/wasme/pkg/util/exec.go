package util

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

func container() string {
	if os.Getenv("WASME_USE_PODMAN") != "" {
		return "podman"
	}
	return "docker"
}

func DockerRun(stdout, stderr io.Writer, stdin io.Reader, image string, runArgs, imageArgs []string) error {
	args := append([]string{"run"}, runArgs...)
	args = append(args, image)
	args = append(args, imageArgs...)
	return ExecCmd(stdout, stderr, stdin, container(), args...)
}

func Docker(stdout, stderr io.Writer, stdin io.Reader, args ...string) error {
	return ExecCmd(stdout, stderr, stdin, container(), args...)
}

func ExecOutput(stdin io.Reader, cmd string, args ...string) (string, error) {
	buf := &bytes.Buffer{}
	err := ExecCmd(buf, buf, stdin, cmd, args...)
	out := strings.TrimSpace(buf.String())
	if err != nil {
		return "", errors.Wrap(err, out)
	}
	return out, nil
}

func ExecCmd(stdout, stderr io.Writer, stdin io.Reader, cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	command.Stderr = stderr
	command.Stdout = stdout
	command.Stdin = stdin
	logrus.Debugf("exec: %v", command.Args)
	return command.Run()
}

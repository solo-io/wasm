package util

import (
	"io"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func DockerRun(stdout, stderr io.Writer, stdin io.Reader, image string, runArgs, imageArgs []string) error {
	args := append([]string{"run"}, runArgs...)
	args = append(args, image)
	args = append(args, imageArgs...)
	return ExecCmd(stdout, stderr, stdin, "docker", args...)
}

func Docker(stdout, stderr io.Writer, stdin io.Reader, args ...string) error {
	return ExecCmd(stdout, stderr, stdin, "docker", args...)
}

func ExecCmd(stdout, stderr io.Writer, stdin io.Reader, cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	command.Stderr = stderr
	command.Stdout = stdout
	command.Stdin = stdin
	logrus.Debugf("exec: %v", command.Args)
	return command.Run()
}

package test

import (
	"github.com/solo-io/wasme/pkg/cmd"
	"strings"
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

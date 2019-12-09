package initialize

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasme/pkg/util"
	"github.com/spf13/cobra"
	"path/filepath"
)

var log = logrus.StandardLogger()

type initOptions struct {
	destDir string
}

func InitCmd() *cobra.Command {
	var opts initOptions
	cmd := &cobra.Command{
		Use:   "init DEST_DIRECTORY",
		Short: "Initialize a source directory for new Envoy WASM Filter.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("invalid number of arguments")
			}
			opts.destDir = args[0]
			return runInit(opts)
		},
	}

	return cmd
}

func runInit(opts initOptions) error {
	destDir, err := filepath.Abs(opts.destDir)
	if err != nil {
		return err
	}

	// currently only supports CPP
	reader := bytes.NewBuffer(cppTarBytes)

	log.Infof("extracting %v bytes to %v", len(cppTarBytes), destDir)

	return util.Untar(destDir, reader)
}

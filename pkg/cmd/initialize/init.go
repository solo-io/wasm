package initialize

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasme/pkg/util"
	"github.com/spf13/cobra"
)

var log = logrus.StandardLogger()

type initOptions struct {
	destDir    string
	filterBase string
}

const (
	filterBaseCppIstio = "cpp-istio" // required for istio 1.4
	filterBaseCpp      = "cpp"
)

// contains map of example projects to the keyword provided by the user
var baseNameToArchive = map[string][]byte{
	filterBaseCppIstio: cppIstio1_4TarBytes,
	filterBaseCpp:      cppTarBytes,
}

var validBases = []string{
	filterBaseCpp,
	filterBaseCppIstio,
}

func InitCmd() *cobra.Command {
	var opts initOptions
	cmd := &cobra.Command{
		Use: "init DEST_DIRECTORY [--base=FILTER_BASE]",
		Short: `Initialize a source directory for new Envoy WASM Filter.

The provided --base will determine the content of the created directory. The default is 
a C++ example filter compatible with the latest Envoy Wasm APIs.

Note that Istio 1.4 uses an older version of the Envoy Wasm APIs and users should 
use --base=cpp-istio to initialize a filter source directory for Istio.
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("invalid number of arguments")
			}
			opts.destDir = args[0]
			return runInit(opts)
		},
	}

	cmd.PersistentFlags().StringVar(&opts.filterBase, "--base", filterBaseCpp,
		fmt.Sprintf("The type of filter to build. Valid filter bases are: %v", validBases))

	return cmd
}

func runInit(opts initOptions) error {
	destDir, err := filepath.Abs(opts.destDir)
	if err != nil {
		return err
	}

	archive, ok := baseNameToArchive[opts.filterBase]
	if !ok {
		return errors.Errorf("%v is not a valid base name. valid names: %v", opts.filterBase, validBases)
	}

	reader := bytes.NewBuffer(archive)

	log.Infof("extracting %v bytes to %v", len(archive), destDir)

	return util.Untar(destDir, reader)
}

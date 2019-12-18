package deploy

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/wasme/pkg/util"
	"github.com/spf13/cobra"
	"path/filepath"
)

var log = logrus.StandardLogger()

const (
	Provider_Gloo  = "gloo"
	Provider_Istio = "istio"
	Provider_Envoy = "envoy"
)

var SupportedProviders = []string{
	Provider_Gloo,
	Provider_Istio,
	Provider_Envoy,
}

type deployOptions struct {
	provider    string
	operation       string
	dryRun          bool
	injectName      string
	injectNamespace string
}

func DeployCmd() *cobra.Command {
	var opts deployOptions
	cmd := &cobra.Command{
		Use:   "deploy DEST_DIRECTORY",
		Short: "Deploy an Envoy WASM Filter to the data plane (Envoy proxies).",
		Long: `Deploys an Envoy WASM Filter to Envoy instances.

Deploy contains a subcommand for each 
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("invalid number of arguments")
			}
			opts.destDir = args[0]
			return runDeploy(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.operation, "operation", "o", "_output/filter.wasm", "path to the output .wasm file. Nonexistent directories will be created.")

	return cmd
}

func runDeploy(opts deployOptions) error {
	destDir, err := filepath.Abs(opts.destDir)
	if err != nil {
		return err
	}

	// currently only supports CPP
	reader := bytes.NewBuffer(cppTarBytes)

	log.Infof("extracting %v bytes to %v", len(cppTarBytes), destDir)

	return util.Untar(destDir, reader)
}

func makeDeployCommand(provider Provider) {

}

func deployGloo() error {

}

func updateGateway(gw *v1.Gateway) error {

}

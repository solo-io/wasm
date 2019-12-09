package pull

import (
	"context"
	"io/ioutil"

	"github.com/sirupsen/logrus"

	"github.com/solo-io/extend-envoy/pkg/cmd/opts"
	"github.com/solo-io/extend-envoy/pkg/pull"
	"github.com/solo-io/extend-envoy/pkg/resolver"
	"github.com/spf13/cobra"
)

type pullOptions struct {
	targetRef          string
	allowedMediaTypes  []string
	allowAllMediaTypes bool
	keepOldFiles       bool
	pathTraversal      bool
	output             string
	verbose            bool

	debug bool

	*opts.GeneralOptions
}

func PullCmd(generalOptions *opts.GeneralOptions) *cobra.Command {
	var opts pullOptions
	opts.GeneralOptions = generalOptions
	cmd := &cobra.Command{
		Use:   "pull <name:tag|name@digest> [-o output-file]",
		Short: "Pull files from remote registry",
		Long: `Pull files from remote registry

Example - Pull only files with the "application/vnd.oci.image.layer.v1.tar" media type (default):
  oras pull localhost:5000/hello:latest

Example - Pull only files with the custom "application/vnd.me.hi" media type:
  oras pull localhost:5000/hello:latest -t application/vnd.me.hi

Example - Pull all files, any media type:
  oras pull localhost:5000/hello:latest -a

Example - Pull files from the insecure registry:
  oras pull localhost:5000/hello:latest --insecure

Example - Pull files from the HTTP registry:
  oras pull localhost:5000/hello:latest --plain-http
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.targetRef = args[0]
			return runPull(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.output, "output", "o", "filter.wasm", "output file")
	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose output")

	cmd.Flags().BoolVarP(&opts.debug, "debug", "d", false, "debug mode")
	return cmd
}

func runPull(opts pullOptions) error {

	ctx := context.Background()
	puller := pull.NewPuller(resolver.NewResolver(opts.Username, opts.Password, opts.Insecure, opts.PlainHTTP, opts.Configs...))

	filter, err := puller.Pull(ctx, opts.targetRef)
	if err != nil {
		return err
	}

	logrus.Printf("Pulled filter image %v", opts.targetRef)

	raw, err := ioutil.ReadAll(filter.Code())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(opts.output, raw, 0644)
}

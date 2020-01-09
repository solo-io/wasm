package pull

import (
	"context"
	"io/ioutil"

	"github.com/sirupsen/logrus"

	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/solo-io/wasme/pkg/pull"
	"github.com/solo-io/wasme/pkg/resolver"
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

func PullCmd(ctx *context.Context, generalOptions *opts.GeneralOptions) *cobra.Command {
	var opts pullOptions
	opts.GeneralOptions = generalOptions
	cmd := &cobra.Command{
		Use:   "pull <name:tag|name@digest> [-o output-file]",
		Short: "Pull wasm filters from remote registry",
		Long: `Pull wasm filters from remote registry
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.targetRef = args[0]
			return runPull(*ctx, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.output, "output", "o", "filter.wasm", "output file")
	return cmd
}

func runPull(ctx context.Context, opts pullOptions) error {
	resolver, _ := resolver.NewResolver(opts.Username, opts.Password, opts.Insecure, opts.PlainHTTP, opts.Configs...)
	var puller pull.ImagePuller = pull.NewPuller(resolver)

	filter, err := puller.PullFilter(ctx, opts.targetRef)
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

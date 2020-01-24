package pull

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasme/pkg/store"

	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/solo-io/wasme/pkg/pull"
	"github.com/solo-io/wasme/pkg/resolver"
	"github.com/spf13/cobra"
)

type pullOptions struct {
	ref        string
	storageDir string

	*opts.GeneralOptions
}

func PullCmd(ctx *context.Context, generalOptions *opts.GeneralOptions) *cobra.Command {
	var opts pullOptions
	opts.GeneralOptions = generalOptions
	cmd := &cobra.Command{
		Use:   "pull <name:tag|name@digest>",
		Short: "Pull wasm filters from remote registry",
		Long: `Pull wasm filters from remote registry
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ref = args[0]
			return runPull(*ctx, opts)
		},
	}

	cmd.Flags().StringVar(&opts.storageDir, "store", "", "Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store")

	return cmd
}

func runPull(ctx context.Context, opts pullOptions) error {
	logrus.Infof("Pulling image %v", opts.ref)

	resolver, _ := resolver.NewResolver(opts.Username, opts.Password, opts.Insecure, opts.PlainHTTP, opts.Configs...)
	var puller pull.ImagePuller = pull.NewPuller(resolver)

	image, err := puller.Pull(ctx, opts.ref)
	if err != nil {
		return err
	}

	desc, err := image.Descriptor()
	if err != nil {
		return err
	}

	if err := store.NewStore(opts.storageDir).Add(ctx, image); err != nil {
		return err
	}

	logrus.Infof("Pulled digest", desc.Digest)

	return nil
}

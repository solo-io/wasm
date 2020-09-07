package pull

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasm/tools/wasme/pkg/store"

	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/opts"
	"github.com/solo-io/wasm/tools/wasme/pkg/pull"
	"github.com/solo-io/wasm/tools/wasme/pkg/resolver"
	"github.com/spf13/cobra"
)

type pullOptions struct {
	ref        string
	storageDir string

	*opts.AuthOptions
}

func PullCmd(ctx *context.Context, loginOptions *opts.AuthOptions) *cobra.Command {
	var opts pullOptions
	opts.AuthOptions = loginOptions
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

	resolver, _ := resolver.NewResolver(opts.Username, opts.Password, opts.Insecure, opts.PlainHTTP, opts.CredentialsFiles...)
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

	logrus.Infof("Image: %v", image.Ref())
	logrus.Infof("Digest: %v", desc.Digest)

	return nil
}

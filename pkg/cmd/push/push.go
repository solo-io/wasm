package push

import (
	"context"
	"fmt"

	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/solo-io/wasme/pkg/push"
	"github.com/solo-io/wasme/pkg/resolver"
	"github.com/spf13/cobra"
)

type pushOptions struct {
	targetRef string
	code      string
	config    string
	verbose   bool

	debug bool

	*opts.GeneralOptions
}

func PushCmd(generalOptions *opts.GeneralOptions) *cobra.Command {
	var opts pushOptions
	opts.GeneralOptions = generalOptions
	cmd := &cobra.Command{
		Use:   "push name[:tag|@digest] code.wasm [config_proto-descriptor-set.proto.bin]",
		Short: "Push wasm filter to remote registry",
		Long: `Push wasm filter to remote registry. E.g.:

wasme push webassemblyhub.io/my/filter:v1 filter.wasm
`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 && len(args) != 3 {
				return fmt.Errorf("invalid number of arguments")
			}
			opts.targetRef = args[0]
			opts.code = args[1]
			if len(args) == 3 {
				opts.config = args[2]
			}
			return runPush(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose output")
	cmd.Flags().BoolVarP(&opts.debug, "debug", "d", false, "debug mode")
	return cmd
}

func runPush(opts pushOptions) error {
	pusher := push.PusherImpl{
		Resolver: resolver.NewResolver(opts.Username, opts.Password, opts.Insecure, opts.PlainHTTP, opts.Configs...),
	}
	return pusher.Push(context.Background(), push.NewLocalFilter(opts.code, opts.config, opts.targetRef))
}

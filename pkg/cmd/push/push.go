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
	targetRef   string
	code        string
	descriptors string
	rootId      string
	verbose     bool

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
				opts.descriptors = args[2]
			}
			return runPush(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.rootId, "root-id", "r", "", "Specify the root_id of the filter to be loaded by Envoy. If not specified, users of this filter will have to specify the --root-id flag to the `wasme deploy` command.")
	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose output")
	cmd.Flags().BoolVarP(&opts.debug, "debug", "d", false, "debug mode")
	return cmd
}

func runPush(opts pushOptions) error {
	resolver, authorizer := resolver.NewResolver(opts.Username, opts.Password, opts.Insecure, opts.PlainHTTP, opts.Configs...)
	pusher := push.PusherImpl{
		Resolver:   resolver,
		Authorizer: authorizer,
	}
	return pusher.Push(context.Background(), push.NewLocalFilter(opts.code, opts.descriptors, opts.targetRef, opts.rootId))
}

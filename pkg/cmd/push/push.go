package push

import (
	"context"
	"fmt"

	"github.com/solo-io/extend-envoy/pkg/cmd/opts"
	"github.com/solo-io/extend-envoy/pkg/push"
	"github.com/solo-io/extend-envoy/pkg/resolver"
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
		Short: "Push wasm to remote registry",
		Long: `Push wasm to remote registry

Example - Push file "hi.txt" with the "application/vnd.oci.image.layer.v1.tar" media type (default):
  oras push localhost:5000/hello:latest hi.txt

Example - Push file "hi.txt" with the custom "application/vnd.me.hi" media type:
  oras push localhost:5000/hello:latest hi.txt:application/vnd.me.hi

Example - Push multiple files with different media types:
  oras push localhost:5000/hello:latest hi.txt:application/vnd.me.hi bye.txt:application/vnd.me.bye

Example - Push file "hi.txt" with the custom manifest config "config.json" of the custom "application/vnd.me.config" media type:
  oras push --manifest-config config.json:application/vnd.me.config localhost:5000/hello:latest hi.txt

Example - Push file to the insecure registry:
  oras push localhost:5000/hello:latest hi.txt --insecure

Example - Push file to the HTTP registry:
  oras push localhost:5000/hello:latest hi.txt --plain-http
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

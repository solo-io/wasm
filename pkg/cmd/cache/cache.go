package cache

import (
	"context"
	"fmt"
	"net/http"

	"github.com/solo-io/extend-envoy/pkg/cmd/opts"
	"github.com/solo-io/extend-envoy/pkg/defaults"
	"github.com/spf13/cobra"
)

type cacheOptions struct {
	targetRefs []string
	code       string
	config     string
	verbose    bool

	debug bool
	port  int

	*opts.GeneralOptions
}

func CacheCmd(generalOptions *opts.GeneralOptions) *cobra.Command {
	var opts cacheOptions
	opts.GeneralOptions = generalOptions
	cmd := &cobra.Command{
		Use:   "cache name[:tag|@digest] ...",
		Short: "Expose images using http and their sha",
		Long: `cache
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("invalid number of arguments")
			}
			opts.targetRefs = args
			return runCache(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose output")
	cmd.Flags().BoolVarP(&opts.debug, "debug", "d", false, "debug mode")
	cmd.Flags().IntVarP(&opts.port, "port", "", 9979, "port")
	return cmd
}

func runCache(opts cacheOptions) error {

	imageCache := defaults.NewDefaultCache()
	for _, image := range opts.targetRefs {
		digest, err := imageCache.Add(context.TODO(), image)
		if err != nil {
			return fmt.Errorf("can't add image")
		}
		fmt.Println("added digest", digest)
	}

	http.ListenAndServe(fmt.Sprintf(":%d", opts.port), imageCache)
	return nil
}

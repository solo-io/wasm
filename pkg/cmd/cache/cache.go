package cache

import (
	"context"
	"fmt"
	"net/http"

	"github.com/solo-io/wasme/pkg/cache"

	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/solo-io/wasme/pkg/defaults"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type cacheOptions struct {
	targetRefs []string
	code       string
	config     string
	verbose    bool

	debug     bool
	port      int
	directory string
	refFile   string

	*opts.AuthOptions
}

func CacheCmd(ctx *context.Context, loginOptions *opts.AuthOptions) *cobra.Command {
	var opts cacheOptions
	opts.AuthOptions = loginOptions
	cmd := &cobra.Command{
		Use:   "cache name[:tag|@digest] ...",
		Short: "Expose images using http and their sha",
		Long: `cache
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 && opts.refFile == "" {
				return fmt.Errorf("invalid number of arguments")
			}
			opts.targetRefs = args
			return runCache(*ctx, opts)
		},
		Hidden: true,
	}

	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "verbose output")
	cmd.Flags().BoolVarP(&opts.debug, "debug", "d", false, "debug mode")
	cmd.Flags().IntVarP(&opts.port, "port", "", 9979, "port")
	cmd.Flags().StringVarP(&opts.directory, "directory", "", "", "directory to write the refs we need to cache")
	cmd.Flags().StringVarP(&opts.refFile, "ref-file", "", "", "file to watch for images we need to cache.")
	return cmd
}

func runCache(ctx context.Context, opts cacheOptions) error {

	imageCache := defaults.NewDefaultCache()
	for _, image := range opts.targetRefs {
		digest, err := imageCache.Add(context.TODO(), image)
		if err != nil {
			return fmt.Errorf("can't add image")
		}
		fmt.Println("added digest", digest)
	}

	errg, ctx := errgroup.WithContext(ctx)

	if 0 != opts.port {
		errg.Go(func() error {
			return http.ListenAndServe(fmt.Sprintf(":%d", opts.port), imageCache)
		})
	}
	if opts.refFile != "" {
		errg.Go(func() error {
			return watchFile(ctx, imageCache, opts.refFile, opts.directory)
		})
	}
	return errg.Wait()
}

func watchFile(ctx context.Context, imageCache cache.Cache, refFile, directory string) error {
	// for each ref in the file, add it to the cache,
	// and if directory is not empty write it t here
	fw := cache.NewLocalImagePuller(
		imageCache,
		refFile,
		directory,
	)

	return fw.WatchFile(ctx)
}

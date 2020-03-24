package cache

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/solo-io/wasme/pkg/cache"

	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/solo-io/wasme/pkg/defaults"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type cacheOptions struct {
	targetRefs []string
	debug      bool
	port       int
	directory  string
	refFile    string
	clearCache bool

	kubeOpts kubeOpts

	*opts.AuthOptions
}

type kubeOpts struct {
	disableKube    bool
	cacheNamespace string
	cacheName      string
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

	cmd.Flags().IntVarP(&opts.port, "port", "", 9979, "port")
	cmd.Flags().StringVarP(&opts.directory, "directory", "", "", "directory to write the refs we need to cache")
	cmd.Flags().StringVarP(&opts.refFile, "ref-file", "", "", "file to watch for images we need to cache.")
	cmd.Flags().BoolVarP(&opts.clearCache, "clear-cache", "", false, "clear any files from the cache dir on boot")
	cmd.Flags().BoolVarP(&opts.kubeOpts.disableKube, "disable-kube", "", true, "disable sending events to kubernetes when images are pulled successfully")
	cmd.Flags().StringVarP(&opts.kubeOpts.cacheNamespace, "cache-ns", "", cache.CacheNamespace, "namespace where the cache is running, if kube integration is enabled")
	cmd.Flags().StringVarP(&opts.kubeOpts.cacheName, "cache-name", "", cache.CacheName, "name of the cache configmap")
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
			return watchFile(ctx, imageCache, opts.refFile, opts.directory, opts.clearCache, opts.kubeOpts)
		})
	}
	return errg.Wait()
}

func watchFile(ctx context.Context, imageCache cache.Cache, refFile, directory string, clearCache bool, kubeOpts kubeOpts) error {

	if clearCache {
		cacheContents, err := ioutil.ReadDir(directory)
		if err != nil {
			return errors.Wrap(err, "reading cache dir")
		}
		for _, file := range cacheContents {
			logrus.Infof("removing cached file %v", file.Name())
			if err := os.RemoveAll(file.Name()); err != nil {
				return err
			}
		}
	}

	var cacheNotifier cache.EventNotifier
	if !kubeOpts.disableKube {
		cfg := config.GetConfigOrDie()
		kube := kubernetes.NewForConfigOrDie(cfg)
		cacheNotifier = cache.NewNotifier(
			kube,
			kubeOpts.cacheNamespace,
			kubeOpts.cacheName,
		)
	}

	// for each ref in the file, add it to the cache,
	// and if directory is not empty write it t here
	fw := cache.NewLocalImagePuller(
		imageCache,
		refFile,
		directory,
		cacheNotifier,
	)

	return fw.WatchFile(ctx)
}

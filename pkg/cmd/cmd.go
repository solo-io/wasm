package cmd

import (
	"os"

	"github.com/solo-io/extend-envoy/pkg/cmd/opts"
	"github.com/solo-io/extend-envoy/pkg/cmd/pull"
	"github.com/solo-io/extend-envoy/pkg/cmd/push"
	"github.com/spf13/cobra"
)

func Run() {
	cmd := &cobra.Command{
		Use: "extend-envoy [command]",
	}
	var opts opts.GeneralOptions
	cmd.AddCommand(push.PushCmd(&opts), pull.PullCmd(&opts))
	cmd.PersistentFlags().StringArrayVarP(&opts.Configs, "config", "c", nil, "auth config path")
	cmd.PersistentFlags().StringVarP(&opts.Username, "username", "u", "", "registry username")
	cmd.PersistentFlags().StringVarP(&opts.Password, "password", "p", "", "registry password")
	cmd.PersistentFlags().BoolVarP(&opts.Insecure, "insecure", "", false, "allow connections to SSL registry without certs")
	cmd.PersistentFlags().BoolVarP(&opts.PlainHTTP, "plain-http", "", false, "use plain http and not https")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

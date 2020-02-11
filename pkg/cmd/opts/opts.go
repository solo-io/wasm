package opts

import "github.com/spf13/pflag"

type AuthOptions struct {
	CredentialsFiles []string
	Username         string
	Password         string
	Insecure         bool
	PlainHTTP        bool

	Verbose bool
	Debug   bool
}

func (opts *AuthOptions) AddToFlags(flags *pflag.FlagSet) {
	flags.StringArrayVarP(&opts.CredentialsFiles, "config", "c", nil, "path to auth config")
	flags.StringVarP(&opts.Username, "username", "u", "", "registry username")
	flags.StringVarP(&opts.Password, "password", "p", "", "registry password")
	flags.BoolVarP(&opts.Insecure, "insecure", "", false, "allow connections to SSL registry without certs")
	flags.BoolVarP(&opts.PlainHTTP, "plain-http", "", false, "use plain http and not https")
	flags.BoolVarP(&opts.Verbose, "verbose", "v", false, "verbose output")
	flags.BoolVarP(&opts.Debug, "debug", "d", false, "debug mode")
}

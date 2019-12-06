package resolver

import (
	"crypto/tls"
	auth "github.com/deislabs/oras/pkg/auth/docker"
	"net/http"
	"strings"

	"github.com/solo-io/extend-envoy/pkg/auth/store"

	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
)

const basicAuthToTokenUser = "basic2token"

func NewResolver(username, password string, insecure bool, plainHTTP bool, configs ...string) remotes.Resolver {

	opts := docker.ResolverOptions{
		PlainHTTP: plainHTTP,
	}

	client := http.DefaultClient
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	opts.Client = client

	if username != "" || password != "" {
		opts.Credentials = func(hostName string) (string, string, error) {
			return username, password, nil
		}
		return docker.NewResolver(opts)
	}
	var dockerCreds func(hostname string) (string, string, error)
	if cli, err := auth.NewClient(configs...); err != nil {
		if authcli, ok := cli.(*auth.Client); ok {
			dockerCreds = authcli.Credential
		}
	}
	resolver, err := cli.Resolver(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading resolver: %v\n", err)
		resolver = docker.NewResolver(opts)

	token, _ := store.GetToken()

	opts.Credentials = func(hostName string) (string, string, error) {

		if token != "" {
			if hostName == "getwasm.io" {
				return basicAuthToTokenUser, token, nil
			}
			if strings.HasPrefix(hostName, "localhost") {
				return basicAuthToTokenUser, token, nil
			}
		}

		if dockerCreds != nil {
			return dockerCreds(hostName)
		}
		return "", "", nil
	}

	return docker.NewResolver(opts)
}

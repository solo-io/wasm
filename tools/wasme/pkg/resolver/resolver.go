package resolver

import (
	"crypto/tls"
	"net/http"

	auth "github.com/deislabs/oras/pkg/auth/docker"

	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
)

func NewResolver(username, password string, insecure bool, plainHTTP bool, configs ...string) (remotes.Resolver, docker.Authorizer) {

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
		return docker.NewResolver(opts), nil
	}
	var dockerCreds func(hostname string) (string, string, error)
	if cli, err := auth.NewClient(configs...); err == nil {
		if authcli, ok := cli.(*auth.Client); ok {
			dockerCreds = authcli.Credential
		}
	}

	credentials := func(hostName string) (string, string, error) {
		if dockerCreds != nil {
			return dockerCreds(hostName)
		}
		return "", "", nil
	}

	opts.Authorizer = docker.NewDockerAuthorizer(
		docker.WithAuthClient(opts.Client),
		docker.WithAuthHeader(opts.Headers),
		docker.WithAuthCreds(credentials),
	)

	return docker.NewResolver(opts), opts.Authorizer
}

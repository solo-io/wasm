package resolver

import (
	"crypto/tls"
	"net/http"
	"strings"

	auth "github.com/deislabs/oras/pkg/auth/docker"

	"github.com/solo-io/wasme/pkg/auth/store"
	"github.com/solo-io/wasme/pkg/consts"

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

	token, _ := store.GetToken()

	opts.Credentials = func(hostName string) (string, string, error) {

		if token != "" {
			if hostName == consts.HubDomain {
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

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
	if cli, err := auth.NewClient(configs...); err != nil {
		if authcli, ok := cli.(*auth.Client); ok {
			dockerCreds = authcli.Credential
		}
	}

	token, err := store.GetToken()
	if err != nil {
		// TODO: log err for pushes
		//logrus.Warnf("Warning: No token found. Make sure to run `wasme login`: %v", err)
	}

	credentials := func(hostName string) (string, string, error) {

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

	opts.Authorizer = docker.NewDockerAuthorizer(docker.WithAuthClient(opts.Client), docker.WithAuthHeader(opts.Headers), docker.WithAuthCreds(credentials))

	return docker.NewResolver(opts), opts.Authorizer
}

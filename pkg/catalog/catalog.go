package catalog

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/v28/github"
	"github.com/solo-io/extend-envoy/pkg/auth/store"
	"golang.org/x/oauth2"
)

func UpdateCatalogItem(ctx context.Context, ref, catalogRepo, itemname, contents string) error {

	token := os.Getenv("GITHUB_API_TOKEN")
	if token == "" {
		var err error
		token, err = store.GetToken()
		if err != nil {
			return err
		}
	}

	// TODO: token
	var httpClient *http.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		httpClient = oauth2.NewClient(ctx, ts)
	}

	client := github.NewClient(httpClient)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return err
	}
	githubUsername := *user.Login
	catalogOwner := "solo-io"
	branch := "master"
	gt := NewGithubTransaction(ctx, client, catalogOwner, catalogRepo, branch, githubUsername, ref)

	steps := []struct {
		name string
		f    func() error
	}{
		{
			"Making sure your fork is available",
			gt.EnsureFork,
		},
		{
			"Making sure a feature branch is available",
			gt.EnsureBranch,
		},
		{
			"Adding your catalog item to the feature branch",
			func() error { return gt.ModifyBranch(itemname, contents) },
		},
		{
			"Makeing sure a PR is open",
			func() error {
				prState, err := gt.CurrentPrState()
				if err != nil {
					return err
				}

				if prState != Pending {
					if err := gt.EnsurePr(); err != nil {
						return err
					}
				}
				return nil
			},
		},
	}

	for _, step := range steps {
		if err := step.f(); err != nil {
			return err
		}
	}

	return nil
}

package catalog

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

func UpdateCatalogItem(ctx context.Context, token, ref, catalogRepo, itemname, contents string) error {
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
	fmt.Println("ensure fork")
	if err := gt.EnsureFork(); err != nil {
		return err
	}
	fmt.Println("ensure branch")
	if err := gt.EnsureBranch(); err != nil {
		return err
	}

	fmt.Println("current pr state")
	prState, err := gt.CurrentPrState()
	if err != nil {
		return err
	}
	if prState == Merged {
		fmt.Println("current pr merged - del bra")
		if err := gt.DeleteBranch(); err != nil {
			return err
		}
		fmt.Println("current pr merged - ens bra")
		if err := gt.EnsureBranch(); err != nil {
			return err
		}
	}
	fmt.Println("modify branch")

	if err := gt.ModifyBranch(itemname, contents); err != nil {
		return err
	}
	if prState != Pending {
		fmt.Println("current pr merged - ens pr")
		if err := gt.EnsurePr(); err != nil {
			return err
		}
	}
	return nil
}

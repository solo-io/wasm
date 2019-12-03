package catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v28/github"
)

type PullRequestState int

const (
	Merged PullRequestState = iota
	None
	Pending
)

type githubTransaction struct {
	ctx    context.Context
	client *github.Client

	originalOwner  string
	originalRepo   string
	originalBranch string
	forkOwner      string
	forkBranch     string
}

func NewGithubTransaction(
	ctx context.Context,
	client *github.Client,
	originalOwner string,
	originalRepo string,
	originalBranch string,
	forkOwner string,
	forkBranch string) *githubTransaction {
	return &githubTransaction{
		ctx:            ctx,
		client:         client,
		originalOwner:  originalOwner,
		originalRepo:   originalRepo,
		originalBranch: originalBranch,
		forkOwner:      forkOwner,
		forkBranch:     forkBranch,
	}

}

func (g *githubTransaction) EnsureFork() error {
	// don't set org as we are forking to the user's account
	rcf := &github.RepositoryCreateForkOptions{}

	// first check that we have the repo
	checkRepo := func() error {
		_, _, err := g.client.Git.GetRef(g.ctx, g.forkOwner, g.originalRepo, "refs/heads/"+g.originalBranch)
		return err
	}

	if checkRepo() == nil {
		return nil
	}

	_, _, err := g.client.Repositories.CreateFork(g.ctx, g.originalOwner, g.originalRepo, rcf)
	if err != nil {
		if _, ok := err.(*github.AcceptedError); ok {
			subctx, cancel := context.WithTimeout(g.ctx, 5*time.Minute)
			defer cancel()
			//loop and wait until the fork is ready
			// try to get the master branch to test that it is indeed ready
			for {

				if checkRepo() == nil {
					return nil
				}
				select {
				case <-subctx.Done():
					return fmt.Errorf("timed out waiting for fork")
				case <-time.After(5 * time.Second):
				}
			}
		}
	}

	// if fork exists we are in good shape
	return err
}

func (g *githubTransaction) EnsureBranch() error {

	masterRefStr := "refs/heads/" + g.originalBranch
	masterRef, _, err := g.client.Git.GetRef(g.ctx, g.forkOwner, g.originalRepo, masterRefStr)
	if err != nil {
		return err
	}

	refstr := "refs/heads/" + g.forkBranch
	_, resp, err := g.client.Git.GetRef(g.ctx, g.forkOwner, g.originalRepo, refstr)
	if err == nil {
		return nil
	}
	if err != nil && resp.StatusCode != 404 {
		return err
	}

	ref := &github.Reference{
		Ref: github.String(refstr),
		Object: &github.GitObject{
			SHA: github.String(masterRef.GetObject().GetSHA()),
		},
	}

	_, _, err = g.client.Git.CreateRef(g.ctx, g.forkOwner, g.originalRepo, ref)
	return err

}

func (g *githubTransaction) CurrentPrState() (PullRequestState, error) {
	// find a pr from githubUsername:branch to catalogRepo:master
	// if all existing PRs are merged, return true
	// if we have no PRs RETURN NONE
	// if we have a pending PR, return pending

	opts := &github.PullRequestListOptions{
		Head:  g.forkOwner + ":" + g.forkBranch,
		State: "all",
	}
	prs, _, err := g.client.PullRequests.List(g.ctx, g.originalOwner, g.originalRepo, opts)
	if err != nil {
		return PullRequestState(0), err
	}

	if len(prs) == 0 {
		return None, nil
	}

	for _, pr := range prs {
		if pr.GetState() == "open" {
			return Pending, nil
		}
	}
	return Merged, nil
}

func (g *githubTransaction) DeleteBranch() error {
	_, err := g.client.Git.DeleteRef(g.ctx, g.forkOwner, g.originalRepo, g.forkBranch)
	return err
}

func (g *githubTransaction) ModifyBranch(file, content string) error {
	opt := &github.RepositoryContentFileOptions{
		Branch:  github.String(g.forkBranch),
		Message: github.String("update catalog"),
		Content: []byte(content),
	}
	_, _, err := g.client.Repositories.CreateFile(g.ctx, g.forkOwner, g.originalRepo, file, opt)
	return err
}

func (g *githubTransaction) EnsurePr() error {
	newPR := &github.NewPullRequest{
		Title: github.String("catalog: add item"),
		Body:  github.String("catalog: add item"),
		Base:  github.String(g.originalBranch),
		Head:  github.String(g.forkOwner + ":" + g.forkBranch),
	}
	_, _, err := g.client.PullRequests.Create(g.ctx, g.originalOwner, g.originalRepo, newPR)
	return err
}

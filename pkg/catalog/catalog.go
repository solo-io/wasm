package catalog

import (
	"context"
	fmt "fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/containerd/containerd/reference"
	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/google/go-github/v28/github"
	"github.com/solo-io/wasme/pkg/auth/store"
	"golang.org/x/oauth2"
)

//go:generate protoc --proto_path=. --go_out=. catalog.proto

const (
	catalogRepo = "wasme"
)

func getName(refspec reference.Spec) string {
	s := strings.Split(refspec.Locator, "/")
	return s[len(s)-1]
}

func getOwner(refspec reference.Spec) string {
	// format should be: host/owner/repo
	s := strings.Split(refspec.Locator, "/")
	if len(s) >= 2 {
		return s[1]
	}
	return ""
}

func refToFolder(refspec reference.Spec) string {
	base := getName(refspec)
	return path.Join("catalog", base, refspec.Object, "spec.yaml")
}

func getContents(refspec reference.Spec) (string, error) {
	m, err := getManifestFromUser(refspec)
	if err != nil {
		return "", err
	}

	var marshal jsonpb.Marshaler

	json, err := marshal.MarshalToString(m)
	if err != nil {
		return "", err
	}
	yamlb, err := yaml.JSONToYAML([]byte(json))
	if err != nil {
		return "", err
	}
	yamls := string(yamlb)
	fmt.Println(m, err, yamls)

	return yamls, fmt.Errorf("no")
}

func UpdateCatalogItem(ctx context.Context, ref string) error {

	refspec, err := reference.Parse(ref)
	if err != nil {
		return err
	}
	itemname := refToFolder(refspec)
	contents, err := getContents(refspec)
	if err != nil {
		return err
	}

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

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
	survey "gopkg.in/AlecAivazis/survey.v1"
)

//go:generate protoc --proto_path=. --go_out=. catalog.proto

const (
	wasmeCatalogRepo = "wasme"
	catalogOwner     = "solo-io"
	branch           = "master"
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

	return yamls, nil
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

	catalogRepo := os.Getenv("CATALOG_REPO")
	if catalogRepo == "" {
		catalogRepo = wasmeCatalogRepo
	}

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

	forkBranch := getName(refspec) + "-" + refspec.Object
	gt := NewGithubTransaction(ctx, client, catalogOwner, catalogRepo, branch, githubUsername, forkBranch)

	steps := []struct {
		desc string
		name string
		f    func() error
	}{
		{
			fmt.Sprintf("Fork github.com/%s/%s", catalogOwner, catalogRepo),
			"Making sure your fork is available",
			gt.EnsureFork,
		},
		{
			fmt.Sprintf("Create a feature branch named %s", forkBranch),
			"Making sure a feature branch is available",
			gt.EnsureBranch,
		},
		{
			"Add your spec to this branch in this location: " + itemname,
			"Adding your catalog item to the feature branch",
			func() error { return gt.ModifyBranch(itemname, contents) },
		},
		{
			fmt.Sprintf("And open a Pull Request against github.com/%s/%s", catalogOwner, catalogRepo),
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

	message := "We are all set to submit your entry to the catalog. Here is your spec:\n" + contents + "\nIn these steps we will:\n"
	for _, step := range steps {
		message += fmt.Sprintln("\t", step.desc)
	}

	prompt := &survey.Confirm{
		Message: message,
	}
	var answer bool
	survey.AskOne(prompt, &answer, nil)
	if !answer {
		fmt.Println("aborted")
		return nil
	}

	for _, step := range steps {
		fmt.Println(step.name)
		if err := step.f(); err != nil {
			return err
		}
	}

	return nil
}

package catalog

import (
	fmt "fmt"

	"github.com/containerd/containerd/reference"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func getManifestFromUser(refspec reference.Spec) (*ModuleSpec, error) {
	getString := func(m string, out *string) error {
		prompt := &survey.Input{
			Message: m,
		}
		if len(*out) != 0 {
			prompt.Default = *out
		}
		return survey.AskOne(prompt, out, nil)
	}
	var ms ModuleSpec

	ms.Name = getName(refspec)
	ms.CreatorName = getOwner(refspec)
	ms.RepositoryUrl = fmt.Sprintf("https://gitub.com/%s/%s", ms.CreatorName, ms.Name)
	ms.ModuleUrl = refspec.String()

	steps := []struct {
		f func() error
	}{
		{
			func() error {
				return getString("Please provide name of the extension", &ms.Name)
			},
		},
		{
			func() error {
				return getString("Please provide short description of the extension", &ms.ShortDescription)
			},
		},
		{
			func() error {
				return getString("Please provide the url to the source code", &ms.RepositoryUrl)
			},
		},
	}

	for _, step := range steps {
		if err := step.f(); err != nil {
			return nil, err
		}
	}

	return &ms, nil
}

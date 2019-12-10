package catalog

import (
	fmt "fmt"

	"github.com/containerd/containerd/reference"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func getManifestFromUser(refspec reference.Spec) (*ExtensionSpec, error) {
	getString := func(m string, out *string) error {
		prompt := &survey.Input{
			Message: m,
		}
		if len(*out) != 0 {
			prompt.Default = *out
		}
		return survey.AskOne(prompt, out, nil)
	}
	getRequiredString := func(m string, out *string) error {
		if *out == "" {
			// no default value, let user know this is required
			m = m + " (required)"
		}
		err := getString(m, out)
		if err == nil && *out == "" {
			return fmt.Errorf("required field not provided")
		}
		return err
	}
	var ms ExtensionSpec

	ms.Name = getName(refspec)
	ms.CreatorName = getOwner(refspec)
	ms.RepositoryUrl = fmt.Sprintf("https://gitub.com/%s/%s", ms.CreatorName, ms.Name)
	ms.ExtensionRef = refspec.String()

	steps := []struct {
		f func() error
	}{
		{
			func() error {
				return getRequiredString("Please provide name of the extension", &ms.Name)
			},
		},
		{
			func() error {
				return getRequiredString("Please provide short description of the extension", &ms.ShortDescription)
			},
		},
		{
			func() error {
				return getString("Please provide long description of the extension", &ms.LongDescription)
			},
		},
		{
			func() error {
				return getString("Please provide the url to the source code", &ms.RepositoryUrl)
			},
		},
		{
			func() error {
				return getString("Please provide the url to the documentation", &ms.DocumentationUrl)
			},
		},
		{
			func() error {
				return getRequiredString("Please provide the name of the extension creator", &ms.CreatorName)
			},
		},
		{
			func() error {
				return getString("Please provide a url for the extension creator", &ms.CreatorUrl)
			},
		},
		{
			func() error {
				return getString("Please provide a url for the extension logo", &ms.LogoUrl)
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

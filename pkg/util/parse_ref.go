package util

import "github.com/docker/distribution/reference"

// splits a ref into the repo and tag
// if tag is empty, returns "latest"
func SplitImageRef(ref string) (string, string, error) {
	named, err := reference.ParseNormalizedNamed(ref)
	if err != nil {
		return "", "", err
	}

	tag := "latest"
	if tagged, isTagged := named.(reference.Tagged); isTagged {
		tag = tagged.Tag()
	}

	return named.Name(), tag, nil
}

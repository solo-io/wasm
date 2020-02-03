package util

import "strings"

// splits a ref into the repo and tag
// if tag is empty, returns "latest"
func SplitImageRef(ref string) (string, string) {
	parts := strings.Split(ref, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return ref, "latest"

}

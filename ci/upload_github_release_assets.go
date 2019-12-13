package main

import (
	"github.com/solo-io/go-utils/githubutils"
)

const buildDir = "_output"
const repoOwner = "solo-io"
const repoName = "wasme"

func main() {
	assets := []githubutils.ReleaseAssetSpec{
		{
			Name:       "wasme-linux-amd64",
			ParentPath: buildDir,
			UploadSHA:  true,
		},
		{
			Name:       "wasme-darwin-amd64",
			ParentPath: buildDir,
			UploadSHA:  true,
		},
		{
			Name:       "wasme-windows-amd64.exe",
			ParentPath: buildDir,
			UploadSHA:  true,
		},
	}
	spec := githubutils.UploadReleaseAssetSpec{
		Owner:             repoOwner,
		Repo:              repoName,
		Assets:            assets,
		SkipAlreadyExists: true,
	}
	githubutils.UploadReleaseAssetCli(&spec)
}

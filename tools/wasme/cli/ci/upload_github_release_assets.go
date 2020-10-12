package main

import (
	"github.com/solo-io/go-utils/githubutils"
)

const buildDir = "_output"
const installDir = "operator/install"
const repoOwner = "solo-io"
const repoName = "wasm"

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
		{
			Name:       "wasme-default.yaml",
			ParentPath: installDir,
			UploadSHA:  true,
		},
		{
			Name:       "wasme.io_v1_crds.yaml",
			ParentPath: installDir + "/wasme/crds",
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

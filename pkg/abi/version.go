package abi

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	PlatformNameIstio = "istio"
	PlatformNameGloo  = "gloo"

	Version13x = "1.3.x"
	Version15x = "1.5.x"
)

// an abi provider
type Platform struct {
	// the name of the mesh supporting
	Name    string
	Version string
}

// a reference to a version of the ABI
type Version struct {
	// an internal logical name for the version
	Name string

	// the git repo URL containing the Envoy ABI
	Repository string

	// the commit SHA of the target ABI version
	Commit string
}

type Versions []Version

// map of ABI Version to the Providers using that version
type Registry map[Version][]Platform

/*
Which platform do you wish to compile against?

Istio 1.4.x [ABI version 6d525c67f39b36cdff9d688697f266c1b55e9cb7] [ ]
Istio 1.5.x [ABI version 541b2c1155fffb15ccde92b8324f3e38f7339ba6] [x]
Gloo  1.3.x [ABI version 541b2c1155fffb15ccde92b8324f3e38f7339ba6] [ ]

*/

func (registry Registry) SelectVersion(platform Platform) (Version, bool) {
	for version, supportedPlatforms := range registry {
		for _, supportedPlatform := range supportedPlatforms {
			if platform == supportedPlatform {
				return version, true
			}
		}
	}
	return Version{}, false
}

// helper check the abi version compatibility
func (registry Registry) ValidateIstioVersion(abiVersion, istioVersion string) error {
	var versionFound bool
	for version, platforms := range registry {
		if version.Name == abiVersion {
			versionFound = true
			for _, platform := range platforms {
				if platform.Name != PlatformNameIstio {
					continue
				}
				match, err := matchVersion(istioVersion, platform.Version)
				if err != nil {
					return err
				}
				if match {
					return nil
				}
			}
		}
	}
	if !versionFound {
		return errors.Errorf("abi version %v not found", abiVersion)
	}
	return errors.Errorf("no versions of istio found which match abi version %v. registered versions: %v", abiVersion, registry)
}

// the default registry of AbiVersions used by Wasme
var (
	Istio15 = Platform{
		Name:    PlatformNameIstio,
		Version: Version15x,
	}
	Gloo13 = Platform{
		Name:    PlatformNameGloo,
		Version: Version13x,
	}

	Version_541b2c1155fffb15ccde92b8324f3e38f7339ba6 = Version{
		Name:       "v0-541b2c1155fffb15ccde92b8324f3e38f7339ba6",
		Repository: "https://github.com/yuval-k/envoy-wasm",
		Commit:     "541b2c1155fffb15ccde92b8324f3e38f7339ba6",
	}
	Version_097b7f2e4cc1fb490cc1943d0d633655ac3c522f = Version{
		Name:       "v0-097b7f2e4cc1fb490cc1943d0d633655ac3c522f",
		Repository: "https://github.com/envoyproxy/envoy-wasm",
		Commit:     "097b7f2e4cc1fb490cc1943d0d633655ac3c522f",
	}

	DefaultRegistry = Registry{
		Version_541b2c1155fffb15ccde92b8324f3e38f7339ba6: {
			Gloo13,
		},
		Version_097b7f2e4cc1fb490cc1943d0d633655ac3c522f: {
			Istio15,
		},
	}
)

// match a real version to an X version, e.g.
// 1.4.2 == 1.4.x
func matchVersion(realVersion, xVersion string) (bool, error) {
	rxp, err := regexp.Compile(strings.ReplaceAll(xVersion, "x", `.*`))
	if err != nil {
		return false, err
	}
	return rxp.MatchString(realVersion), nil
}

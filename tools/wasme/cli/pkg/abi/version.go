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
	Version16x = "1.6.x"
	Version17x = "1.7.x"
	Version18x = "1.8.x"
	Version19x = "1.9.x"
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
func (registry Registry) ValidateIstioVersion(abiVersions []string, istioVersion string) error {
	var versionFound bool
	for version, platforms := range registry {
		for _, abiVersion := range abiVersions {
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
	}
	if !versionFound {
		return errors.Errorf("abi versions %v not found", abiVersions)
	}
	return errors.Errorf("no versions of istio found which support abi versions %v. registered versions: %v", abiVersions, registry)
}

// the default registry of AbiVersions used by Wasme
var (
	Istio15 = Platform{
		Name:    PlatformNameIstio,
		Version: Version15x,
	}
	Istio16 = Platform{
		Name:    PlatformNameIstio,
		Version: Version16x,
	}
	Istio17 = Platform{
		Name:    PlatformNameIstio,
		Version: Version17x,
	}
	Istio18 = Platform{
		Name:    PlatformNameIstio,
		Version: Version18x,
	}
	Istio19 = Platform{
		Name:    PlatformNameIstio,
		Version: Version19x,
	}
	Gloo13 = Platform{
		Name:    PlatformNameGloo,
		Version: Version13x,
	}
	Gloo15 = Platform{
		Name:    PlatformNameGloo,
		Version: Version15x,
	}
	Gloo16 = Platform{
		Name:    PlatformNameGloo,
		Version: Version16x,
	}

	Version_097b7f2e4cc1fb490cc1943d0d633655ac3c522f = Version{
		// December 12 2019
		Name:       "v0-097b7f2e4cc1fb490cc1943d0d633655ac3c522f",
		Repository: "https://github.com/envoyproxy/envoy-wasm",
		Commit:     "097b7f2e4cc1fb490cc1943d0d633655ac3c522f",
	}
	Version_edc016b1fa5adca3ebd3d7020eaed0ad7b8814ca = Version{
		// July 7th 2020
		Name:       "v0-edc016b1fa5adca3ebd3d7020eaed0ad7b8814ca",
		Repository: "https://github.com/envoyproxy/envoy-wasm",
		Commit:     "edc016b1fa5adca3ebd3d7020eaed0ad7b8814ca",
	}

	// Precursor to ABI v0_2_0
	Version_4689a30309abf31aee9ae36e73d34b1bb182685f = Version{
		// August 4th 2020
		Name:       "v0-4689a30309abf31aee9ae36e73d34b1bb182685f",
		Repository: "https://github.com/envoyproxy/envoy-wasm",
		Commit:     "4689a30309abf31aee9ae36e73d34b1bb182685f",
	}

	Version_0_2_1 = Version{
		// Oct 23rd 2020
		Name:       "v0.2.1",
		Repository: "https://github.com/envoyproxy/envoy",
		Commit:     "2758575b9a02f935245fe8a3b08af0a8c14994dc",
	}

	DefaultRegistry = Registry{
		Version_edc016b1fa5adca3ebd3d7020eaed0ad7b8814ca: {
			Gloo15,
		},
		Version_097b7f2e4cc1fb490cc1943d0d633655ac3c522f: {
			Istio15,
			Istio16,
		},
		Version_4689a30309abf31aee9ae36e73d34b1bb182685f: {
			Istio17,
			Istio18,
			Istio19,
		},
		Version_0_2_1: {
			Gloo16,
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

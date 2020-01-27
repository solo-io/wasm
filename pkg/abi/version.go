package abi

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

func (r Registry) SelectVersion(platform Platform) (Version, bool) {
	for version, supportedPlatforms := range r {
		for _, supportedPlatform := range supportedPlatforms {
			if platform == supportedPlatform {
				return version, true
			}
		}
	}
	return Version{}, false
}

/*
Which platform do you wish to compile against?

Istio 1.4.x [ABI version 6d525c67f39b36cdff9d688697f266c1b55e9cb7] [ ]
Istio 1.5.x [ABI version 541b2c1155fffb15ccde92b8324f3e38f7339ba6] [x]
Gloo  1.3.x [ABI version 541b2c1155fffb15ccde92b8324f3e38f7339ba6] [ ]

*/

const (
	PlatformNameIstio = "istio"
	PlatformNameGloo  = "gloo"

	Version13x = "1.3.x"
	Version14x = "1.4.x"
	Version15x = "1.5.x"
)

// the default registry of AbiVersions used by Wasme
var (
	Istio14 = Platform{
		Name:    PlatformNameIstio,
		Version: Version14x,
	}
	Istio15 = Platform{
		Name:    PlatformNameIstio,
		Version: Version15x,
	}
	Gloo13 = Platform{
		Name:    PlatformNameGloo,
		Version: Version13x,
	}

	VersionIstio14 = Version{
		Name:       "v0-6d525c67f39b36cdff9d688697f266c1b55e9cb7",
		Repository: "https://github.com/istio/envoy",
		Commit:     "6d525c67f39b36cdff9d688697f266c1b55e9cb7",
	}
	VersionIstio15 = Version{
		Name:       "v0-541b2c1155fffb15ccde92b8324f3e38f7339ba6",
		Repository: "https://github.com/yuval-k/envoy-wasm",
		Commit:     "541b2c1155fffb15ccde92b8324f3e38f7339ba6",
	}

	DefaultRegistry = Registry{
		VersionIstio14: {
			Istio14,
		},
		VersionIstio15: {
			Gloo13,
			Istio15,
		},
	}
)

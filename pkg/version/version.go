package version

// This will be set by the linker on release builds
var (
	Version    string
	DevVersion = "dev"
)

func init() {
	if Version == "" {
		Version = DevVersion
	}
}

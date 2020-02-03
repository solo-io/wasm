package initialize

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/manifoldco/promptui"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasme/pkg/abi"
	"github.com/solo-io/wasme/pkg/util"
	"github.com/spf13/cobra"
)

const (
	languageCpp = "cpp"
)

// map supported languages and abi versions to the archive for those sources
var languageVersionArchives = map[string]map[abi.Version][]byte{
	languageCpp: {
		abi.VersionIstio14: cppIstio1_4TarBytes,
		abi.VersionIstio15: cppTarBytes,
	},
}

// map of language name to description
var supportedLanguages = []string{
	languageCpp,
}

func selectSourceArchive(language string, platform abi.Platform) ([]byte, error) {
	version, ok := abi.DefaultRegistry.SelectVersion(platform)
	if !ok {
		return nil, errors.Errorf("no version available for platform %+v", platform)
	}

	languageArchives, ok := languageVersionArchives[language]
	if !ok {
		return nil, errors.Errorf("%v is not a supported language. available: %v", language, supportedLanguages)
	}

	versionedArchive, ok := languageArchives[version]
	if !ok {
		return nil, errors.Errorf("%v is not a supported platform for %v. available: %v", platform, language, supportedPlatforms(language))
	}

	return versionedArchive, nil
}

// list the platforms supported by the language
func supportedPlatforms(language string) []abi.Platform {
	var platforms []abi.Platform
	for version := range languageVersionArchives[language] {
		for _, platform := range abi.DefaultRegistry[version] {
			platforms = append(platforms, platform)
		}
	}
	sort.SliceStable(platforms, func(i, j int) bool {
		return platforms[i].Name < platforms[j].Name && platforms[i].Version < platforms[j].Version
	})
	return platforms
}

var log = logrus.StandardLogger()

type initOptions struct {
	destDir       string
	language      string
	platform      abi.Platform
	disablePrompt bool
}

func InitCmd() *cobra.Command {
	var opts initOptions
	cmd := &cobra.Command{
		Use: "init DEST_DIRECTORY [--language=FILTER_LANGUAGE] [--platform=TARGET_PLATFORM] [--platform-version=TARGET_PLATFORM_VERSION]",
		Short: fmt.Sprintf(`Initialize a project directory for a new Envoy WASM Filter.

The provided --language flag will determine the programming language used for the new filter. The default is 
C++.

The provided --platform flag will determine the target platform used for the new filter. This is important to 
ensure compatibility between the filter and the 

If --language, --platform, or --platform-version are not provided, the CLI will present an interactive prompt. Disable the prompt with --disable-prompt

`),
		Args: cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if opts.language == "" {
				opts.language, err = getLanguageInteractive()
				if err != nil {
					return err
				}
			}
			if opts.platform.Name == "" || opts.platform.Version == "" {
				opts.platform, err = getPlatformInteractive(opts.language)
				if err != nil {
					return err
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("invalid number of arguments")
			}
			opts.destDir = args[0]
			return runInit(opts)
		},
	}

	cmd.PersistentFlags().StringVar(&opts.language, "language", "",
		fmt.Sprintf("The programming language with which to create the filter. Supported languages are: %v", supportedLanguages))

	cmd.PersistentFlags().StringVar(&opts.platform.Name, "platform", "",
		fmt.Sprintf("The name of the target platform against which to build. Supported platforms are: %v", []string{"gloo", "istio"}))

	cmd.PersistentFlags().StringVar(&opts.platform.Version, "platform-version", "",
		fmt.Sprintf("The version of the target platform against which to build. Supported Istio versions are: %v. Supported Gloo versions are: %v", []string{abi.Version14x, abi.Version15x}, []string{abi.Version13x}))

	cmd.PersistentFlags().BoolVar(&opts.disablePrompt, "disable-prompt", false,
		"Disable the interactive prompt if a required parameter is not passed. If set to true and one of the required flags is not provided, wasme CLI will return an error.")

	return cmd
}

func runInit(opts initOptions) error {
	destDir, err := filepath.Abs(opts.destDir)
	if err != nil {
		return err
	}

	archive, err := selectSourceArchive(opts.language, opts.platform)
	if err != nil {
		return err
	}

	reader := bytes.NewBuffer(archive)

	log.Infof("extracting %v bytes to %v", len(archive), destDir)

	return util.Untar(destDir, reader)
}

func getValueInteractive(message string, options interface{}) (string, error) {
	prompt := promptui.Select{
		Label: message,
		Items: options,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func getLanguageInteractive() (string, error) {
	return getValueInteractive(
		"What language do you wish to use for the filter",
		supportedLanguages,
	)
}

func getPlatformInteractive(language string) (abi.Platform, error) {
	var platformOptions []string
	selectablePlatforms := map[string]abi.Platform{}

	for _, platform := range supportedPlatforms(language) {
		key := platform.Name + " " + platform.Version
		selectablePlatforms[key] = platform
		platformOptions = append(platformOptions, key)
	}

	platformKey, err := getValueInteractive(
		"With which platform do you wish to use the filter?",
		platformOptions,
	)
	if err != nil {
		return abi.Platform{}, err
	}
	return selectablePlatforms[platformKey], nil
}

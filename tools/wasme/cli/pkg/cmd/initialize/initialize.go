package initialize

/*
TODO (ilackarms): Devise a better strategy for adding support for new abi versions, platforms, languages and filter bases.
The current steps required to add a new language:
1. Add it to the examples dir
2. Make sure the runtime-config.json is present and contains necessary fields.
3. run `make generated-code` to regen the 2gobytes archives with the new example
4. Add the new example as a new filterBase to the availableBases map below
  - it may be required to add additional ABI Versions to the Registry in abi.DefaultRegistry
*/

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/abi"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasm/tools/wasme/pkg/util"
	"github.com/spf13/cobra"
)

const (
	languageCpp            = "cpp"
	languageRust           = "rust"
	languageAssemblyScript = "assemblyscript"
)

// a filterBase is a starter filter we generate from the example filters
type filterBase struct {
	// the set of versions compatible with a given filter base
	compatiblePlatforms compatiblePlatforms

	// baked-in  bytes, generated with 2-go-array
	archiveBytes []byte
}

// set of compatible abi versions for a single base
type compatiblePlatforms []abi.Platform

// return a sorted list of each name+version of the platforms
func (c compatiblePlatforms) Keys() []string {
	var platformNames []string
	for _, platform := range c {
		platformNames = append(platformNames, platform.Name+":"+platform.Version)
	}
	return platformNames
}

// returns true if c is a superset of c2
func (c compatiblePlatforms) IsSupersetOf(theirs compatiblePlatforms) bool {
	for _, theirPlatform := range theirs {
		var weContainTheirPlatform bool
		for _, ourPlatform := range c {
			if ourPlatform == theirPlatform {
				weContainTheirPlatform = true
				break
			}
		}
		if !weContainTheirPlatform {
			return false
		}
	}
	return true
}

// map supported languages and abi versions to the archive for those sources
var availableBases = map[string][]filterBase{
	languageCpp: {
		{
			// cpp for istio 1.5
			compatiblePlatforms: compatiblePlatforms{
				abi.Istio15,
				abi.Istio16,
			},
			archiveBytes: cppIstio1_5TarBytes,
		},
		{
			// cpp for istio 1.7
			compatiblePlatforms: compatiblePlatforms{
				abi.Istio17,
			},
			archiveBytes: cppIstio1_7TarBytes,
		},
		{
			compatiblePlatforms: compatiblePlatforms{
				abi.Gloo15,
			},
			archiveBytes: cppTarBytes,
		},
	},
	languageAssemblyScript: {
		{
			compatiblePlatforms: compatiblePlatforms{
				abi.Gloo13,
				abi.Istio15,
				abi.Istio16,
			},
			archiveBytes: assemblyscriptTarBytes,
		},
	},
	languageRust: {
		{
			// rust for istio 1.7
			compatiblePlatforms: compatiblePlatforms{
				abi.Istio15,
				abi.Istio16,
				abi.Istio17,
			},
			archiveBytes: rustIstio1_7TarBytes,
		},
	},
}

// map of language name to description
var supportedLanguages = []string{
	languageCpp,
	languageRust,
	languageAssemblyScript,
}

var log = logrus.StandardLogger()

type initOptions struct {
	destDir       string
	language      string
	platform      abi.Platform
	disablePrompt bool

	// set by PreRun
	compatiblePlatforms compatiblePlatforms
}

func InitCmd() *cobra.Command {
	var opts initOptions
	cmd := &cobra.Command{
		Use:   "init DEST_DIRECTORY [--language=FILTER_LANGUAGE] [--platform=TARGET_PLATFORM] [--platform-version=TARGET_PLATFORM_VERSION]",
		Short: fmt.Sprintf(`Initialize a project directory for a new Envoy WASM Filter.`),
		Long: `
The provided --language flag will determine the programming language used for the new filter. The default is 
C++.

The provided --platform flag will determine the target platform used for the new filter. This is important to 
ensure compatibility between the filter and the 

If --language, --platform, or --platform-version are not provided, the CLI will present an interactive prompt. Disable the prompt with --disable-prompt

`,
		Args: cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !opts.disablePrompt {
				var err error
				if opts.language == "" {
					opts.language, err = selectLanguageInteractive()
					if err != nil {
						return err
					}
				}
				if opts.platform.Name == "" || opts.platform.Version == "" {
					opts.compatiblePlatforms, err = selectCompatiblePlatformsInteractive(opts.language)
					if err != nil {
						return err
					}
				} else {
					opts.compatiblePlatforms = compatiblePlatforms{opts.platform}
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
		fmt.Sprintf("The version of the target platform against which to build. Supported Istio versions are: %v. Supported Gloo versions are: %v", []string{abi.Version15x, abi.Version16x}, []string{abi.Version13x, abi.Version15x}))

	cmd.PersistentFlags().BoolVar(&opts.disablePrompt, "disable-prompt", false,
		"Disable the interactive prompt if a required parameter is not passed. If set to true and one of the required flags is not provided, wasme CLI will return an error.")

	return cmd
}

func runInit(opts initOptions) error {
	destDir, err := filepath.Abs(opts.destDir)
	if err != nil {
		return err
	}

	base, err := getFilterBase(opts.language, opts.compatiblePlatforms)
	if err != nil {
		return err
	}

	reader := bytes.NewBuffer(base.archiveBytes)

	log.Infof("extracting %v bytes to %v", len(base.archiveBytes), destDir)

	return util.Untar(destDir, reader)
}

func selectLanguageInteractive() (string, error) {
	return selectValueInteractive(
		"What language do you wish to use for the filter",
		supportedLanguages,
	)
}

// the user selects a set of supported platforms from a filter base
func selectCompatiblePlatformsInteractive(language string) (compatiblePlatforms, error) {
	bases, ok := availableBases[language]
	if !ok {
		return nil, errors.Errorf("%v is not a supported language. available: %v", language, supportedLanguages)
	}

	var baseOptions []string
	selectableBases := map[string]*filterBase{}

	for _, base := range bases {
		base := base
		key := strings.Join(base.compatiblePlatforms.Keys(), ", ")
		selectableBases[key] = &base
		baseOptions = append(baseOptions, key)
	}

	baseKey, err := selectValueInteractive(
		"With which platforms do you wish to use the filter?",
		baseOptions,
	)
	if err != nil {
		return nil, err
	}
	return selectableBases[baseKey].compatiblePlatforms, nil
}

func selectValueInteractive(message string, options interface{}) (string, error) {
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

func getFilterBase(language string, platforms compatiblePlatforms) (*filterBase, error) {
	bases, ok := availableBases[language]
	if !ok {
		return nil, errors.Errorf("%v is not a supported language. available: %v", language, supportedLanguages)
	}

	for _, base := range bases {
		if base.compatiblePlatforms.IsSupersetOf(platforms) {
			return &base, nil
		}
	}

	return nil, errors.Errorf("no filter base found for language %v is not a supported language. available: %v", language, supportedLanguages)
}

package auth

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"

	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	"github.com/solo-io/wasme/pkg/consts"
	"github.com/solo-io/wasme/pkg/defaults"
)

func SaveCredentials(username, password, serverAddress, path string) error {
	if serverAddress == "" {
		serverAddress = consts.HubDomain
	}
	if path == "" {
		path = defaults.WasmeCredentialsFile
	}
	cfg := configfile.New(path)

	cfg.AuthConfigs[serverAddress] = types.AuthConfig{
		Username: username,
		Password: password,
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil{
		return err
	}

	credsFile, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "creating file")
	}
	defer credsFile.Close()

	if err := cfg.SaveToWriter(credsFile); err != nil {
		return err
	}

	return nil
}

package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func getTokenFile() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir = filepath.Join(dir, "wasm-hub")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "token"), nil
}

func SaveToken(t string) error {
	f, err := getTokenFile()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f, []byte(t), 0600)
}

func GetToken() (string, error) {

	f, err := getTokenFile()
	if err != nil {
		return "", err
	}

	token, err := ioutil.ReadFile(f)
	return strings.TrimSpace(string(token)), err
}

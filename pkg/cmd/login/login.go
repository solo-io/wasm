package login

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasme/pkg/consts"
	"github.com/solo-io/wasme/pkg/defaults"

	"github.com/solo-io/wasme/pkg/auth"
	"github.com/spf13/cobra"
)

type loginOptions struct {
	credentialsFile string
	username        string
	password        string
	serverAddress   string
	usePlaintext    bool
}

func LoginCmd() *cobra.Command {
	var opts loginOptions
	cmd := &cobra.Command{
		Use:   "login [-s SERVER_ADDRESS] -u USERNAME -p PASSWORD ",
		Short: "Log in in so you can push images to the remote server.",
		Long: `
Caches credentials for image pushes in the provided credentials-file (defaults to $HOME/.wasme/credentials.json).

Provide -s=SERVER_ADDRESS to provide login credentials for a registry other than webassemblyhub.io.

`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(opts)
		},
	}

	cmd.Flags().StringVar(&opts.credentialsFile, "credentials-file", "", "write to this credentials file. defaults to $HOME/.wasme/credentials.json")
	cmd.Flags().StringVarP(&opts.username, "username", "u", "", "login username")
	cmd.Flags().StringVarP(&opts.password, "password", "p", "", "login password")
	cmd.Flags().StringVarP(&opts.serverAddress, "server", "s", consts.HubDomain, "the address of the remote registry to which to authenticate")
	cmd.Flags().BoolVar(&opts.usePlaintext, "plaintext", false, "use plaintext to connect to the remote registry (HTTP) rather than HTTPS")

	return cmd
}

func runLogin(opts loginOptions) error {
	if opts.credentialsFile == "" {
		opts.credentialsFile = defaults.WasmeCredentialsFile
	}
	if opts.username == "" {
		return errors.Errorf("must specify username")
	}
	if opts.password == "" {
		return errors.Errorf("must specify password")
	}
	usr, err := getCurrentUser(opts.username, opts.password, opts.serverAddress, opts.usePlaintext)
	if err != nil {
		return err
	}
	logrus.Infof("Successfully logged in as %v (%v)", opts.username, usr.Realname)
	if err := auth.SaveCredentials(opts.username, opts.password, opts.serverAddress, opts.credentialsFile); err != nil {
		return err
	}
	logrus.Infof("stored credentials in %v", opts.credentialsFile)
	return nil
}

func getCurrentUser(username, password, registryAddr string, usePlaintext bool) (*user, error) {
	scheme := "https"
	if usePlaintext {
		scheme = "http"
	}
	req, err := http.NewRequest("GET", scheme+"://"+registryAddr+"/api/users/current", nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("non-200 status code: %vL: %v", res.StatusCode, res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var usr user
	if err := json.Unmarshal(b, &usr); err != nil {
		return nil, err
	}

	return &usr, nil
}

type user struct {
	UserID          int       `json:"user_id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	PasswordVersion string    `json:"password_version"`
	Realname        string    `json:"realname"`
	Comment         string    `json:"comment"`
	Deleted         bool      `json:"deleted"`
	RoleName        string    `json:"role_name"`
	RoleID          int       `json:"role_id"`
	SysadminFlag    bool      `json:"sysadmin_flag"`
	AdminRoleInAuth bool      `json:"admin_role_in_auth"`
	ResetUUID       string    `json:"reset_uuid"`
	CreationTime    time.Time `json:"creation_time"`
	UpdateTime      time.Time `json:"update_time"`
}

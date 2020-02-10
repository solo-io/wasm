package opts

type AuthOptions struct {
	CredentialsFiles []string
	Username         string
	Password         string
	Insecure         bool
	PlainHTTP        bool

	Verbose bool
	Debug   bool
}

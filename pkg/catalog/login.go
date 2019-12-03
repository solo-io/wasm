package catalog

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func HubEndpoint() *url.URL {

	endpoint := os.Getenv("HUB_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://getwasm.io/"
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}
	return u
}

func ResolveHubEndpoint(path string) *url.URL {

	return HubEndpoint().ResolveReference(&url.URL{Path: path})
}

func Login() error {
	// start http server on a random port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}

	authUrl := ResolveHubEndpoint("/authorize")

	port := listener.Addr().(*net.TCPAddr).Port
	currentQuery := authUrl.Query()
	currentQuery.Add("port", strconv.Itoa(port))
	authUrl.RawQuery = currentQuery.Encode()

	fmt.Println("Using port:", port)
	fmt.Println("go to: ", authUrl.String())
	handler, accessTokenChan := NewHandler()
	go http.Serve(listener, handler)
	accessToken := <-accessTokenChan

	fmt.Println("success ! ", accessToken)
	return nil

}

func NewHandler() (http.Handler, <-chan string) {
	c := make(chan string, 1)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		if code != "" {
			// redirect to the wasm endpoint so that it can exchange the code
			codeUrl := ResolveHubEndpoint("/github/callback")
			codeUrl.RawQuery = r.URL.RawQuery
			http.Redirect(w, r, codeUrl.String(), http.StatusSeeOther)
			return
		}

		token := r.FormValue("token")
		select {
		case c <- token:
		default:
		}
		w.Write([]byte("auth success!"))
	})

	return h, c
}

package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/solo-io/extend-envoy/pkg/auth/store"
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

func Login(ctx context.Context) error {
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
	urlString := authUrl.String()

	if err := openBrowser(urlString); err != nil {
		fmt.Println("Cannot launch browser. Please open this url in your browser: ", urlString)
	} else {
		fmt.Println("Opening browser for login. If the browser did not open for you, please go to: ", urlString)
	}

	handler, accessTokenChan := NewHandler()
	go http.Serve(listener, handler)
	select {
	case accessToken := <-accessTokenChan:
		if err := store.SaveToken(accessToken); err != nil {
			panic(err)
		}
		fmt.Println("success ! you are now authenticated")
	case <-ctx.Done():
		return ctx.Err()
	}

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
		w.Write([]byte(`
		<html><head></head><body>
		auth success! You can now close this window.
		</body></html>
		`))
	})

	return h, c
}

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	err := cmd.Start()
	go cmd.Wait()

	return err
}

package list

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/solo-io/wasme/pkg/consts"

	"github.com/solo-io/wasme/pkg/util"

	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasme/pkg/store"
	"github.com/spf13/cobra"
)

type listOpts struct {
	published  bool
	wide       bool
	showDir    bool
	server     string
	search     string
	storageDir string
}

func ListCmd() *cobra.Command {
	var opts listOpts
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Envoy WASM Filters stored locally or published to webassemblyhub.io.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.published, "published", "", false, "Set to true to list images that have been published to webassemblyhub.io. Defaults to listing image stored in local image cache.")
	cmd.Flags().BoolVarP(&opts.wide, "wide", "w", false, "Set to true to list images with their full tag length.")
	cmd.Flags().BoolVarP(&opts.showDir, "show-dir", "d", false, "Set to true to show the local directories for images. Does not apply to published images.")
	cmd.Flags().StringVarP(&opts.server, "server", "s", consts.HubDomain, "If using --published, read images from this remote registry.")
	cmd.Flags().StringVarP(&opts.search, "search", "", "", "If using --published, search images from the remote registry. If unset, `wasme list --published` will return the top repositories which are accessed the most.")
	cmd.Flags().StringVar(&opts.storageDir, "store", "", "Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store. Ignored if using --published")

	return cmd
}

func runList(opts listOpts) error {
	var images []image
	if opts.published {
		i, err := getPublishedImages(opts.server, opts.search)
		if err != nil {
			return err
		}
		images = i
	} else {
		i, err := getLocalImages(opts.storageDir)
		if err != nil {
			return err
		}
		images = i
	}

	sort.Slice(images, func(i, j int) bool {
		if images[i].name < images[j].name {
			return true
		}
		if images[i].name > images[j].name {
			return false
		}
		return images[i].updated.Before(images[j].updated)
	})

	showDir := !opts.published && opts.showDir

	buf := os.Stdout

	// create a new tabwriter
	w := new(tabwriter.Writer)

	w.Init(buf, 0, 0, 0, ' ', 0)

	line := "NAME \tTAG \tSIZE \tSHA \tUPDATED\n"
	if showDir {
		line = "NAME \tTAG \tSIZE \tSHA \tUPDATED\tDIRECTORY\n"
	}
	fmt.Fprintf(w, line)
	for _, image := range images {
		image.Write(w, opts.wide, showDir)
	}
	w.Flush()
	return nil
}

type image struct {
	name      string
	sum       string
	updated   time.Time
	tag       string
	sizeBytes int64

	// only applicable for local images
	dir string
}

func (i image) Write(w io.Writer, wide, showDir bool) {
	sum := i.sum
	if len(sum) > 8 {
		sum = strings.TrimPrefix(sum, "sha256:")[:8]
	}
	tag := i.tag
	if !wide && len(tag) > 32 {
		tag = strings.TrimPrefix(tag, "sha256:")[:32] + "..."
	}

	args := []interface{}{
		i.name, tag, byteCountSI(i.sizeBytes), sum, i.updated.Format(time.RFC822),
	}
	line := "%v \t%v \t%v \t%v \t%v\n"

	if showDir {
		args = append(args, i.dir)
		line = "%v \t%v \t%v \t%v \t%v \t%v\n"
	}

	fmt.Fprintf(w, line, args...)
}

func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func getLocalImages(storageDir string) ([]image, error) {
	imageStore := store.NewStore(storageDir)

	storedImages, err := imageStore.List()
	if err != nil {
		return nil, err
	}

	var images []image
	for _, img := range storedImages {
		name, tag, err := util.SplitImageRef(img.Ref())
		if err != nil {
			logrus.Errorf("failed parsing image ref %v: %v", img.Ref(), err)
			continue
		}

		descriptor, err := img.Descriptor()
		if err != nil {
			return nil, err
		}

		dir, err := imageStore.Dir(img.Ref())
		if err != nil {
			logrus.Errorf("failed getting image %v dir: %v", img.Ref(), err)
			continue
		}

		imageInfo, err := os.Stat(dir)
		if err != nil {
			return nil, err
		}

		images = append(images, image{
			name:      name,
			sum:       descriptor.Digest.String(),
			updated:   imageInfo.ModTime(),
			tag:       tag,
			sizeBytes: descriptor.Size,
			dir:       dir,
		})
	}

	return images, nil
}

func getPublishedImages(serverAddress, searchQuery string) ([]image, error) {
	var repos []repository
	if searchQuery != "" {
		r, err := searchRepos(serverAddress, searchQuery)
		if err != nil {
			return nil, err
		}
		repos = r
	} else {
		r, err := getTopRepos(serverAddress)
		if err != nil {
			return nil, err
		}
		repos = r
	}
	var images []image
	for _, repo := range repos {
		tags, err := getTags(serverAddress, repo.Name)
		if err != nil {
			return nil, err
		}
		for _, tag := range tags {
			images = append(images, image{
				name:      serverAddress + "/" + repo.Name,
				sum:       tag.Digest,
				updated:   tag.PushTime,
				tag:       tag.Name,
				sizeBytes: tag.Size,
			})
		}
	}
	return images, nil
}

func getTopRepos(serverAddress string) ([]repository, error) {
	var repos []repository
	err := getJson(serverAddress, fmt.Sprintf("/api/repositories/top"), &repos)
	return repos, err
}

func searchRepos(serverAddress, query string) ([]repository, error) {
	var searchRes searchResult
	err := getJson(serverAddress, fmt.Sprintf("/api/search?q=%v", query), &searchRes)

	var repos []repository
	for _, repo := range searchRes.Repository {
		repos = append(repos, repository{Name: repo.RepositoryName})
	}

	return repos, err
}

func getTags(serverAddress, repo string) ([]tag, error) {
	var tags []tag
	err := getJson(serverAddress, fmt.Sprintf("/api/repositories/%v/tags?detail=true", repo), &tags)
	return tags, err
}

func getJson(serverAddress, path string, into interface{}) error {
	res, err := http.Get(fmt.Sprintf("https://" + serverAddress + path))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &into); err != nil {
		return err
	}
	return nil
}

type repository struct {
	Name string `json:"name"`
}

type tag struct {
	Digest        string      `json:"digest"`
	Name          string      `json:"name"`
	Size          int64       `json:"size"`
	Architecture  string      `json:"architecture"`
	Os            string      `json:"os"`
	OsVersion     string      `json:"os.version"`
	DockerVersion string      `json:"docker_version"`
	Author        string      `json:"author"`
	Created       time.Time   `json:"created"`
	Config        interface{} `json:"config"`
	Immutable     bool        `json:"immutable"`
	Annotations   struct {
		ModuleWasmRuntimeAbiVersion string `json:"module.wasm.runtime/abi_version"`
		ModuleWasmRuntimeType       string `json:"module.wasm.runtime/type"`
	} `json:"annotations"`
	Signature interface{}   `json:"signature"`
	Labels    []interface{} `json:"labels"`
	PushTime  time.Time     `json:"push_time"`
	PullTime  time.Time     `json:"pull_time"`
}

type searchResult struct {
	Repository []struct {
		ProjectID      int    `json:"project_id"`
		ProjectName    string `json:"project_name"`
		ProjectPublic  bool   `json:"project_public"`
		PullCount      int    `json:"pull_count"`
		RepositoryName string `json:"repository_name"`
		TagsCount      int    `json:"tags_count"`
	} `json:"repository"`
}

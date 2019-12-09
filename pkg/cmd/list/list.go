package list

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

var log = logrus.StandardLogger()

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Envoy WASM Filters published to getwasm.io.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList()
		},
	}

	return cmd
}

func runList() error {
	images, err := getImages()
	if err != nil {
		return err
	}

	buf := os.Stdout

	// create a new tabwriter
	w := new(tabwriter.Writer)

	w.Init(buf, 0, 0, 0, ' ', 0)

	fmt.Fprintf(w, "NAME \tSHA \tUPDATED \tSIZE \tTAGS\n")
	for _, image := range images {
		image.Write(w)
	}
	w.Flush()
	return nil
}

type image struct {
	name      string
	sum       string
	updated   time.Time
	tags      []string
	sizeBytes int64
}

func (i image) Write(w io.Writer) {
	for idx, tag := range i.tags {
		if idx == 0 {
			fmt.Fprintf(w, "%v \t%v \t%v \t%v \t%v\n", i.name, i.sum, i.updated.Format(time.RFC822), byteCountSI(i.sizeBytes), tag)
		} else {
			fmt.Fprintf(w, "  \t  \t  \t  \t%v\n", tag)
		}
	}
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

func getImages() ([]image, error) {
	root, err := getTagInfo("")
	if err != nil {
		return nil, err
	}
	repos := root.Child
	var images []image
	for _, repo := range repos {
		repoInfo, err := getTagInfo(repo)
		if err != nil {
			logrus.Warnf("failed to get repo info for %v, skipping", repo)
			continue
		}
		repoImages := repoInfo.Child

		for _, img := range repoImages {
			imgName := repo + "/" + img
			imgInfo, err := getTagInfo(imgName)
			if err != nil {
				logrus.Warnf("failed to get image info for %v, skipping: %v", repo, err)
				continue
			}
			for sha, manifest := range imgInfo.Manifest {
				if len(sha) > 8 {
					sha = strings.TrimPrefix(sha, "sha256:")[:8]
				}
				size, err := strconv.Atoi(manifest.ImageSizeBytes)
				if err != nil {
					return nil, err
				}
				updated, err := strconv.Atoi(manifest.TimeUploadedMs)
				if err != nil {
					return nil, err
				}
				images = append(images, image{
					name:      imgName,
					sum:       sha,
					updated:   time.Unix(int64(updated), 0),
					tags:      manifest.Tag,
					sizeBytes: int64(size),
				})
			}
		}
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].name < images[j].name
	})

	return images, nil
}

func getTagInfo(repo string) (*tagInfo, error) {
	if repo != "" {
		repo = strings.TrimSuffix(repo, "/") + "/"
	}
	res, err := http.Get(fmt.Sprintf("https://getwasm.io/v2/%vtags/list", repo))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var info tagInfo
	if err := json.Unmarshal(b, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

type tagInfo struct {
	Child    []string            `json:"child"`
	Manifest map[string]manifest `json:"manifest"`
	Name     string              `json:"name"`
	Tags     []string            `json:"tags"`
}

type manifest struct {
	ImageSizeBytes string     `json:"imageSizeBytes"`
	LayerID        string    `json:"layerId"`
	MediaType      string    `json:"mediaType"`
	Tag            []string  `json:"tag"`
	TimeCreatedMs  string `json:"timeCreatedMs"`
	TimeUploadedMs string `json:"timeUploadedMs"`
}

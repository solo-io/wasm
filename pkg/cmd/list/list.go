package list

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/pkg/errors"

	"github.com/solo-io/wasme/pkg/util"

	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasme/pkg/consts"
	"github.com/solo-io/wasme/pkg/store"
	"github.com/spf13/cobra"
)

type listOpts struct {
	published  bool
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
	cmd.Flags().StringVar(&opts.storageDir, "store", "", "Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store. Ignored if using --published")

	return cmd
}

func runList(opts listOpts) error {
	var images []image
	if opts.published {
		i, err := getPublishedImages()
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

	buf := os.Stdout

	// create a new tabwriter
	w := new(tabwriter.Writer)

	w.Init(buf, 0, 0, 0, ' ', 0)

	fmt.Fprintf(w, "NAME \tTAG \tSIZE \tSHA \tUPDATED\n")
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
	tag       string
	sizeBytes int64
}

func (i image) Write(w io.Writer) {
	sum := i.sum
	if len(sum) > 8 {
		sum = strings.TrimPrefix(sum, "sha256:")[:8]
	}
	tag := i.tag
	if len(tag) > 8 {
		tag = strings.TrimPrefix(tag, "sha256:")[:8]
	}

	fmt.Fprintf(w, "%v \t%v \t%v \t%v \t%v\n", i.name, tag, byteCountSI(i.sizeBytes), sum, i.updated.Format(time.RFC822))
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
		})
	}

	sort.Slice(images, func(i, j int) bool {
		if images[i].name < images[j].name {
			return true
		}
		return images[i].updated.Before(images[j].updated)
	})

	return images, nil
}

func getPublishedImages() ([]image, error) {
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
				image, err := parsePublishedImage(imgName, sha, manifest)
				if err != nil {
					// this is a debug line as old images didn't require a tag
					logrus.Debugf("failed to parse info for %v, skipping: %v", imgName, err)
					continue
				}
				images = append(images, image)
			}
		}
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].name < images[j].name
	})

	return images, nil
}

func parsePublishedImage(name, sha string, manifest manifest) (image, error) {
	size, err := strconv.Atoi(manifest.ImageSizeBytes)
	if err != nil {
		return image{}, err
	}
	if len(manifest.Tag) < 1 {
		return image{}, errors.Errorf("invalid manifest, missing tag")
	}
	tag := manifest.Tag[0]
	updated, err := strconv.Atoi(manifest.TimeUploadedMs)
	if err != nil {
		return image{}, err
	}
	return image{
		name:      name,
		sum:       sha,
		updated:   time.Unix(int64(updated/1000), 0),
		tag:       tag,
		sizeBytes: int64(size),
	}, nil
}

func getTagInfo(repo string) (*tagInfo, error) {
	if repo != "" {
		repo = strings.TrimSuffix(repo, "/") + "/"
	}
	res, err := http.Get(fmt.Sprintf("https://"+consts.HubDomain+"/v2/%vtags/list", repo))
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
	ImageSizeBytes string   `json:"imageSizeBytes"`
	LayerID        string   `json:"layerId"`
	MediaType      string   `json:"mediaType"`
	Tag            []string `json:"tag"`
	TimeCreatedMs  string   `json:"timeCreatedMs"`
	TimeUploadedMs string   `json:"timeUploadedMs"`
}

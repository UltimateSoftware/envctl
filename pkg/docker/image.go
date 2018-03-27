package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/alecthomas/template"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/google/uuid"
)

var dockerfileTpl = `FROM {{ .BaseImage }}

VOLUME ["{{ .Mount }}"]

WORKDIR "{{ .Mount }}"

ENTRYPOINT ["{{ .Shell }}"]
`

// BuildImage will build an image based on the passed in ImageConfig. It returns
// the name of the built image, as <cfg.BaseName:UUID>, or an error.
//
// BuildImage blocks until the image build has finished and the API is done
// streaming the output back.
func (c *Client) BuildImage(cfg ImageConfig) (string, error) {
	dockerfile, err := buildDockerfile(cfg)
	if err != nil {
		return "", err
	}

	buildContext, err := getBuildContext(dockerfile)
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("%v:%v", cfg.BaseName, uuid.New().String())
	bldopts := types.ImageBuildOptions{
		Tags: []string{name},
	}

	resp, err := c.ImageBuild(context.Background(), buildContext, bldopts)
	if err != nil {
		return "", err
	}

	// the read MUST happen, if not the program will continue without waiting
	// for the build to complete
	io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	return name, nil
}

// RemoveImage will remove all images with the given `name`. It returns an error
// if something went wrong. An error can be returned with the removal process
// incomplete (if, for example, there was an error removing an image from the
// middle of a list of images).
func (c *Client) RemoveImage(name string) error {
	args := filters.NewArgs()
	args.Add("reference", name)

	rmopts := types.ImageRemoveOptions{
		PruneChildren: true,
		Force:         true,
	}

	imgs, err := c.ImageList(context.Background(), types.ImageListOptions{
		Filters: args,
	})
	if err != nil {
		return err
	}

	for _, img := range imgs {
		_, err := c.ImageRemove(context.Background(), img.ID, rmopts)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildDockerfile(cfg ImageConfig) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	tpl, err := template.New("Dockerfile").Parse(dockerfileTpl)
	if err != nil {
		return &bytes.Buffer{}, err
	}

	err = tpl.Execute(buf, &cfg)
	if err != nil {
		return &bytes.Buffer{}, err
	}

	return buf, nil
}

func getBuildContext(raw *bytes.Buffer) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	wr := tar.NewWriter(buf)

	hdr := &tar.Header{
		Name: "Dockerfile",
		Mode: 0600,
		Size: int64(raw.Len()),
	}

	if err := wr.WriteHeader(hdr); err != nil {
		return &bytes.Buffer{}, err
	}

	if _, err := wr.Write(raw.Bytes()); err != nil {
		return &bytes.Buffer{}, err
	}

	padlen := 512 - (buf.Len() % 512)
	padding := make([]byte, padlen)

	padded := bytes.NewBuffer(append(buf.Bytes(), padding...))

	return padded, nil
}

package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/docker/go-connections/nat"

	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/alecthomas/template"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/google/uuid"
)

func (c *Controller) Create(m container.Metadata) (container.Metadata, error) {
	img, err := c.buildImage(m)
	if err != nil {
		return container.Metadata{}, err
	}

	m.ImageID = img

	cpmap := getContainerPortMappings(m.Ports)
	hpmap := getHostPortMappings(m.Ports)

	ccfg := &docker.Config{
		User:         m.User,
		Tty:          true,
		Image:        m.ImageID,
		OpenStdin:    true,
		Env:          m.Envs,
		ExposedPorts: cpmap,
	}

	hcfg := &docker.HostConfig{
		Binds:        make([]string, 1),
		PortBindings: hpmap,
	}

	hcfg.Binds[0] = m.Mount.String()

	ncfg := &network.NetworkingConfig{}

	cnt, err := c.client.ContainerCreate(
		context.Background(),
		ccfg,
		hcfg,
		ncfg,
		m.BaseName,
	)
	if err != nil {
		return container.Metadata{}, err
	}

	m.ID = cnt.ID
	return m, nil
}

var dockerfileTpl = `FROM {{ .BaseImage }}
	VOLUME ["{{ .Mount.Destination }}"]
	WORKDIR "{{ .Mount.Destination }}"
	ENTRYPOINT ["{{ .Shell }}"]`

// buildImage will build an image based on the passed in ImageConfig. It returns
// the name of the built image, as <cfg.BaseName:UUID>, or an error.
//
// buildImage blocks until the image build has finished and the API is done
// streaming the output back.
func (c *Controller) buildImage(m container.Metadata) (string, error) {
	dockerfile, err := buildDockerfile(m)
	if err != nil {
		return "", err
	}

	buildContext, err := getBuildContext(dockerfile)
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("%v:%v", m.BaseName, uuid.New().String())
	bldopts := types.ImageBuildOptions{
		Tags:    []string{name},
		NoCache: m.NoCache,
	}

	resp, err := c.client.ImageBuild(context.Background(), buildContext, bldopts)
	if err != nil {
		return "", err
	}

	// the read MUST happen, if not the program will continue without waiting
	// for the build to complete
	io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	return name, nil
}

func buildDockerfile(m container.Metadata) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	tpl, err := template.New("Dockerfile").Parse(dockerfileTpl)
	if err != nil {
		return &bytes.Buffer{}, err
	}

	err = tpl.Execute(buf, &m)
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

func getContainerPortMappings(l3map map[string][]int) map[nat.Port]struct{} {
	mappings := map[nat.Port]struct{}{}

	for layer, plist := range l3map {
		for _, p := range plist {
			pstr := nat.Port(fmt.Sprintf("%v/%v", p, layer))
			mappings[pstr] = struct{}{}
		}
	}

	return mappings
}

func getHostPortMappings(l3map map[string][]int) nat.PortMap {
	mappings := nat.PortMap{}

	for layer, plist := range l3map {
		for _, p := range plist {
			pstr := nat.Port(fmt.Sprintf("%v/%v", p, layer))

			bind := nat.PortBinding{
				HostPort: string(pstr),
			}

			mappings[pstr] = []nat.PortBinding{bind}
		}
	}

	return mappings
}

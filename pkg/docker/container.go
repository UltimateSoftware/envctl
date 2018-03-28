package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// CreateContainer creates a container based on the passed in config. The
// container isn't started or attached, just created. It's *not* set to auto
// remove, so it will need to be cleaned up when finished. Any mounts specified
// in the config are used to create bound volumes. Any warnings are swallowed.
// It returns the container's ID or an error if there was an issue.
func (c *Client) CreateContainer(cfg ContainerConfig) (string, error) {
	ccfg := &container.Config{
		User:      "root",
		Tty:       true,
		Image:     cfg.ImageName,
		OpenStdin: true,
		Env:       cfg.Env,
	}

	hcfg := &container.HostConfig{
		Binds: make([]string, len(cfg.Mounts)),
	}

	ncfg := &network.NetworkingConfig{}

	for i, mount := range cfg.Mounts {
		hcfg.Binds[i] = mount.String()
	}

	cnt, err := c.ContainerCreate(
		context.Background(),
		ccfg,
		hcfg,
		ncfg,
		cfg.Name,
	)
	if err != nil {
		return "", err
	}

	return cnt.ID, nil
}

// RemoveContainer removes the container with the given `id`.
func (c *Client) RemoveContainer(id string) error {
	return c.ContainerRemove(
		context.Background(),
		id,
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		},
	)
}

// AttachContainer attaches the terminal session of the currently running
// program to the container interactively.
func (c *Client) AttachContainer(id string) error {
	restoreStdout, restoreStdin, err := c.MakeRawTerminal()
	if err != nil {
		return err
	}

	acfg := types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	}

	resp, err := c.ContainerAttach(context.Background(), id, acfg)
	defer resp.Close()
	if err != nil {
		return err
	}

	errchan := make(chan error)
	donechan := make(chan struct{})

	err = c.ContainerStart(
		context.Background(),
		id,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return err
	}

	go func() {
		_, err := io.Copy(c.stdout.stream, resp.Reader)
		if err != nil {
			errchan <- err
			restoreStdout()
			return
		}

		restoreStdout()
		donechan <- struct{}{}
	}()

	go func() {
		_, err := io.Copy(resp.Conn, c.stdin.stream)
		if err != nil {
			errchan <- err
			restoreStdin()
			return
		}

		restoreStdin()
		resp.CloseWrite()
	}()

	// Depending on the underlying image's entrypoint, there could be cases
	// where there's no command prompt. This could trick the user into thinking
	// that the process is hung, when in fact there just hasn't been anything
	// to write to stdout.
	fmt.Fprintf(
		c.stdout.stream,
		"If you don't see a command prompt, try pressing enter.\r\n",
	)

	select {
	case err = <-errchan:
		if err != nil {
			return err
		}
	case <-donechan:
	}

	return nil
}

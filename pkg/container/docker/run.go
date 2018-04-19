package docker

import (
	"context"
	"io"
	"os"

	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/docker/docker/api/types"
)

// Run runs the given command array on the container with the given metadata.
func (c *Controller) Run(m container.Metadata, cmd []string) (err error) {
	ctx, cancel := context.WithCancel(context.Background())

	c.mirrorContainerTTY(m.ID)

	err = c.client.ContainerStart(
		ctx,
		m.ID,
		types.ContainerStartOptions{},
	)
	if err != nil {
		cancel()
		return err
	}

	cfg := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
		Detach:       false,
		Tty:          true,
	}

	resp, err := c.client.ContainerExecCreate(ctx, m.ID, cfg)
	if err != nil {
		cancel()
		return err
	}

	hijacked, err := c.client.ContainerExecAttach(ctx, resp.ID, cfg)
	if err != nil {
		cancel()
		return err
	}

	errchan := make(chan error)
	donechan := make(chan struct{})
	go func(
		cancel context.CancelFunc,
		hijacked types.HijackedResponse,
		stdout *os.File,
	) {
		_, err = io.Copy(stdout, hijacked.Reader)
		if err != nil {
			cancel()
			errchan <- err
		}

		donechan <- struct{}{}
	}(cancel, hijacked, c.stdout.stream)

	err = c.client.ContainerExecStart(
		ctx,
		resp.ID,
		types.ExecStartCheck{},
	)
	if err != nil {
		cancel()
		return err
	}

	select {
	case err := <-errchan:
		return err
	case <-donechan:
	}

	return nil
}

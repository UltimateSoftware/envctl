package docker

import (
	"context"

	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/docker/docker/api/types"
)

// Run runs the given command array on the container with the given metadata.
func (c *Controller) Run(m container.Metadata, cmd []string) error {
	err := c.client.ContainerStart(
		context.Background(),
		m.ID,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return err
	}

	resp, err := c.client.ContainerExecCreate(
		context.Background(),
		m.ID,
		types.ExecConfig{
			AttachStderr: true,
			AttachStdout: true,
			Cmd:          cmd,
		})
	if err != nil {
		return err
	}

	return c.client.ContainerExecStart(
		context.Background(),
		resp.ID,
		types.ExecStartCheck{},
	)
}

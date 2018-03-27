package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

// RunOnContainer runs the given command array on the container specified by
// cntID.
func (c *Client) RunOnContainer(cmd []string, cntID string) error {
	err := c.ContainerStart(
		context.Background(),
		cntID,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return err
	}

	resp, err := c.ContainerExecCreate(
		context.Background(),
		cntID,
		types.ExecConfig{
			AttachStderr: true,
			AttachStdout: true,
			Cmd:          cmd,
		})
	if err != nil {
		return err
	}

	return c.ContainerExecStart(
		context.Background(),
		resp.ID,
		types.ExecStartCheck{},
	)
}

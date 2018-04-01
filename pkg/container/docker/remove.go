package docker

import (
	"context"

	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/docker/docker/api/types"
)

// RemoveContainer removes the container with the given `id`.
func (c *Controller) RemoveContainer(m container.Metadata) error {
	return c.client.ContainerRemove(
		context.Background(),
		m.ID,
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		},
	)
}

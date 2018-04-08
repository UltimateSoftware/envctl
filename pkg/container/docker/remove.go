package docker

import (
	"context"

	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Remove removes the container with the given metadata.
func (c *Controller) Remove(m container.Metadata) error {
	err := c.removeImage(m.ImageID)
	if err != nil {
		return err
	}

	return c.client.ContainerRemove(
		context.Background(),
		m.ID,
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		},
	)
}

func (c *Controller) removeImage(name string) error {
	args := filters.NewArgs()
	args.Add("reference", name)

	rmopts := types.ImageRemoveOptions{
		PruneChildren: true,
		Force:         true,
	}

	lsopts := types.ImageListOptions{
		Filters: args,
	}

	imgs, err := c.client.ImageList(context.Background(), lsopts)
	if err != nil {
		return err
	}

	for _, img := range imgs {
		_, err := c.client.ImageRemove(context.Background(), img.ID, rmopts)
		if err != nil {
			return err
		}
	}

	return nil
}

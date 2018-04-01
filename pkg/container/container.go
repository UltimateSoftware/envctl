package container

import "fmt"

// Metadata is what's returned by the container functions. It contains
// everything that a consumer of this package needs to know about containers
// being managed.
type Metadata struct {
	ID        string
	ImageID   string
	BaseName  string
	BaseImage string
	Shell     string
	Mounts    []Mount
	WorkDir   string
	Envs      []string
}

// Mount is directory on the host paired with a volume mount point.
type Mount struct {
	Source      string
	Destination string
}

// Controller can control containers. This includes allowing consumers to
// attach to the container.
type Controller interface {
	Create(Metadata) (Metadata, error)
	Remove(Metadata) error
	Attach(Metadata) error
	Run(Metadata, []string) error
}

func (b Mount) String() string {
	return fmt.Sprintf("%v:%v", b.Source, b.Destination)
}

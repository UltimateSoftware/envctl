package container

// Metadata is what's returned by the container functions. It contains
// everything that a consumer of this package needs to know about containers
// being managed.
type Metadata struct {
	ID        string
	BaseName  string
	BaseImage string
	Shell     string
	Mount     string
	WorkDir   string
	Envs      []string
}

// Controller can control containers. This includes allowing consumers to
// attach to the container.
type Controller interface {
	Create(Metadata) (Metadata, error)
	Remove(Metadata) error
	Attach(Metadata) error
	Run(Metadata, []string) error
}

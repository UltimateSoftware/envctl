package container

import "fmt"

// Metadata is what's returned by the container functions. It contains
// everything that a consumer of this package needs to know about containers
// being managed.
type Metadata struct {
	ID        string   `json:"id"`
	ImageID   string   `json:"image_id"`
	BaseName  string   `json:"base_name"`
	BaseImage string   `json:"base_image"`
	Shell     string   `json:"shell"`
	Mount     Mount    `json:"mount"`
	Envs      []string `json:"envs"`
	NoCache   bool     `json:"no_cache"`
	User      string   `json:"user"`
}

// Mount is directory on the host paired with a volume mount point.
type Mount struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

// Controller can control containers. This includes allowing consumers to
// attach to the container.
type Controller interface {
	Create(Metadata) (Metadata, error)
	Remove(Metadata) error
	Attach(Metadata) error
	Run(Metadata, []string) error
}

func (m Mount) String() string {
	return fmt.Sprintf("%v:%v", m.Source, m.Destination)
}

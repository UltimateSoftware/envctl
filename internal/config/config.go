package config

// Opts is what tells envctl what the environment looks like.
type Opts struct {
	Image string `yaml:"image"`
	// The default for this field is true, so `nil`` needs to be discernable
	// from the default `false` value.
	CacheImage *bool `yaml:"cache_image,omitempty"`

	User      string            `yaml:"user"`
	Shell     string            `yaml:"shell"`
	Mount     string            `yaml:"mount,omitempty"`
	Variables map[string]string `yaml:"variables,omitempty"`
	Bootstrap []string          `yaml:"bootstrap,omitempty"`

	// Exposing the host network isn't a cross-platform solution, so the
	// upfront requirement is to expose any ports that the user needs. The ports
	// are to be mapped directly from container to host so that whatever is
	// exposed in the container is the port that's accessed on the host.
	Ports L3Ports `yaml:"ports,omitempty"`
}

// Loader is anything that can load a configuration file.
type Loader interface {
	Load() (Opts, error)
}

// L3Ports are mappings between a layer 3 protocol like TCP and a port number.
type L3Ports map[string][]int

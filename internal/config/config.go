package config

// Opts is what tells envctl what the environment looks like.
type Opts struct {
	Image string `yaml:"image"`
	// The default for this field is true, so `nil`` needs to be discernable
	// from the default `false` value.
	CacheImage *bool `yaml:"cache_image,omitempty"`

	Shell     string            `yaml:"shell"`
	Mount     string            `yaml:"mount,omitempty"`
	Variables map[string]string `yaml:"variables,omitempty"`
	Bootstrap []string          `yaml:"bootstrap,omitempty"`
}

// Loader is anything that can load a configuration file.
type Loader interface {
	Load() (Opts, error)
}

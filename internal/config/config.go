package config

// Opts is what tells envctl what the environment looks like.
type Opts struct {
	Image     string            `yaml:"image"`
	Shell     string            `yaml:"shell"`
	Mount     string            `yaml:"mount,omitempty"`
	Variables map[string]string `yaml:"variables,omitempty"`
	Bootstrap []string          `yaml:"bootstrap,omitempty"`
}

// Loader is anything that can load a configuration file.
type Loader interface {
	Load() (Opts, error)
}

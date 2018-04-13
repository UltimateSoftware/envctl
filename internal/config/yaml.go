package config

import (
	"errors"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// YAML is a Loader for a YAML configuration file.
type YAML struct {
	Path string
}

// Load returns a new `Opts`` by reading the YAML file. If an error
// happens along the way it returns it along with a zeroed `Opts`. If
// something is missing that should be there, it'll return an error.
func (c YAML) Load() (Opts, error) {
	f, err := ioutil.ReadFile(c.Path)
	if err != nil {
		return Opts{}, err
	}

	var cfg Opts
	err = yaml.UnmarshalStrict(f, &cfg)
	if err != nil {
		return Opts{}, err
	}

	if cfg.Image == "" {
		return Opts{}, errors.New("missing image")
	}

	if cfg.Shell == "" {
		return Opts{}, errors.New("missing shell")
	}

	return cfg, nil
}

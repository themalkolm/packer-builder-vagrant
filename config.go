package main

import (
	"fmt"
	"errors"

	"github.com/mitchellh/packer/template/interpolate"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/builder/virtualbox/ovf"
)

type Config struct {
	Name     string `mapstructure:"box_name"`
	Version  string `mapstructure:"box_version"`
	Provider string `mapstructure:"box_provider"`
	BoxFile  string `mapstructure:"box_file"`

	Config   map[string]interface{} `mapstructure:"config"`

	builder  packer.Builder
	ctx      interpolate.Context
}

func NewConfig(raws ...interface{}) (*Config, []string, error) {
	c := new(Config)
	err := config.Decode(c, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &c.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
			},
		},
	}, raws...)
	if err != nil {
		return nil, nil, err
	}

	var errs *packer.MultiError
	if c.Name == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("box_name is required"))
	}
	if c.Provider == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("box_provider is required"))
	}
	if c.BoxFile == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("box_file is required"))
	}
	if c.Version == "" {
		c.Version = "0"
	}

	if _, ok := c.Config["source_path"]; !ok {
		sourcePath, err := findBoxFile(c.Name, c.Version, c.Provider, c.BoxFile)
		if err != nil {
			errs = packer.MultiErrorAppend(errs, err)
		} else {
			c.Config["source_path"] = sourcePath
		}
	}

	var warnings []string = nil
	switch c.Provider {
	case "virtualbox":
		c.builder = &ovf.Builder{}
		warnings, err = c.builder.Prepare(c.Config)
		if err != nil {
			errs = packer.MultiErrorAppend(errs, err)
		}
	default:
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("unsupported provider: %s", c.Provider))
	}

	// Check for any errors.
	if errs != nil && len(errs.Errors) > 0 {
		return nil, warnings, errs
	}

	return c, warnings, nil
}


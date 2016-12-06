package main

import (
	"errors"
	"fmt"

	"github.com/mitchellh/packer/command"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
)

type Config struct {
	Name          string `mapstructure:"box_name"`
	URL           string `mapstructure:"box_url"`
	Version       string `mapstructure:"box_version"`
	Provider      string `mapstructure:"box_provider"`
	BoxFile       string `mapstructure:"box_file"`

	BuilderConfig map[string]interface{} `mapstructure:"builder"`

	builder       packer.Builder
	ctx           interpolate.Context
}

func NewConfig(raws ...interface{}) (*Config, []string, error) {
	c := new(Config)
	err := config.Decode(c, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &c.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
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

	if _, ok := c.BuilderConfig["source_path"]; !ok {
		// TODO(themalkolm): Probably we should not download boxes in ctor. So.. should we move builder creation to
		// Builder::Run then?
		sourcePath, err := fetchBoxFile(c.URL, c.Name, c.Version, c.Provider, c.BoxFile)
		if err != nil {
			errs = packer.MultiErrorAppend(errs, err)
		} else {
			c.BuilderConfig["source_path"] = sourcePath
		}
	}

	builderType, err := c.builderType()
	if err != nil {
		errs = packer.MultiErrorAppend(errs, err)
		return c, nil, nil
	}

	var warnings []string = nil
	if b, ok := command.Builders[builderType]; ok {
		c.builder = b

		warnings, err = c.builder.Prepare(c.BuilderConfig)
		if err != nil {
			errs = packer.MultiErrorAppend(errs, err)
		}
	} else {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("unsupported builder type: %s", builderType))
	}

	// Check for any errors.
	if errs != nil && len(errs.Errors) > 0 {
		return nil, warnings, errs
	}

	return c, warnings, nil
}

func (c *Config) builderType() (string, error) {
	raw, ok := c.BuilderConfig["type"]
	if !ok {
		return "", errors.New("invalid builder config, missing type")
	}

	t, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("invalid builder type value: %#v", t)
	}

	return t, nil
}

package main

import (
	"errors"
	"fmt"

	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig  `mapstructure:",squash"`

	Name          string `mapstructure:"box_name"`
	URL           string `mapstructure:"box_url"`
	Version       string `mapstructure:"box_version"`
	Provider      string `mapstructure:"box_provider"`
	BoxFile       string `mapstructure:"box_file"`

	BuilderConfig map[string]interface{} `mapstructure:"builder"`

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

	// Check for any errors.
	if errs != nil && len(errs.Errors) > 0 {
		return nil, nil, errs
	}

	return c, nil, nil
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

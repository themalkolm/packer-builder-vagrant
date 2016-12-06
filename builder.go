package main

import (
	"fmt"

	"github.com/mitchellh/packer/command"
	"github.com/mitchellh/packer/packer"
)

type Builder struct {
	config  *Config
	builder packer.Builder
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	c, warnings, errs := NewConfig(raws...)
	if errs != nil {
		return warnings, errs
	}
	b.config = c

	return warnings, nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	c := b.config
	if _, ok := c.BuilderConfig["source_path"]; !ok {
		sourcePath, err := fetchBoxFile(c.URL, c.Name, c.Version, c.Provider, c.BoxFile)
		if err != nil {
			return nil, err
		}
		c.BuilderConfig["source_path"] = sourcePath
	}

	builderType, err := c.builderType()
	if err != nil {
		return nil, err
	}

	builder, found := command.Builders[builderType];
	if !found {
		return nil, fmt.Errorf("unsupported builder type: %s", builderType)
	}

	b.builder = builder
	_, err = b.builder.Prepare(c.BuilderConfig)
	if err != nil {
		return nil, err
	}
	return b.builder.Run(ui, hook, cache)
}

func (b *Builder) Cancel() {
	b.builder.Cancel()
}

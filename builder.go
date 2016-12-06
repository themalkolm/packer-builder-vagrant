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
	ui.Message("(vagrant) Builder source_path ...")
	c := b.config
	if _, ok := c.BuilderConfig["source_path"]; !ok {
		sourcePath, err := fetchBoxFile(c.URL, c.Name, c.Version, c.Provider, c.BoxFile)
		if err != nil {
			return nil, err
		}
		c.BuilderConfig["source_path"] = sourcePath
	}
	ui.Message(fmt.Sprintf("(vagrant) Builder source_path: %s", c.BuilderConfig["source_path"]))

	ui.Message("(vagrant) Builder type ...")
	builderType, err := c.builderType()
	if err != nil {
		return nil, err
	}
	ui.Message(fmt.Sprintf("(vagrant) Builder type: %s", builderType))

	ui.Message("(vagrant) Builder ...")
	builder, found := command.Builders[builderType];
	if !found {
		return nil, fmt.Errorf("unsupported builder type: %s", builderType)
	}
	ui.Message("(vagrant) Builder: OK")

	ui.Message("(vagrant) Builder prepare ...")
	b.builder = builder
	warnings, err := b.builder.Prepare(c.BuilderConfig)
	if err != nil {
		return nil, err
	}
	if warnings != nil && len(warnings) > 0 {
		for _, w := range warnings {
			ui.Message(fmt.Sprintf("(vagrant) WARNING: %s", w))
		}
	}
	ui.Message("(vagrant) Builder prepare: OK")

	ui.Message("(vagrant) Builder run ...")
	a, err := b.builder.Run(ui, hook, cache)
	if err != nil {
		return nil, err
	}
	ui.Message("(vagrant) Builder run: OK")
	return a, nil
}

func (b *Builder) Cancel() {
	b.builder.Cancel()
}

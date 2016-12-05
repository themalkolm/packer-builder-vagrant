package main

import (
	"github.com/mitchellh/packer/packer"
)

type Builder struct {
	config *Config
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
	return b.config.builder.Run(ui, hook, cache)
}

func (b *Builder) Cancel() {
	b.config.builder.Cancel()
}

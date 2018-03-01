package main

import (
	"errors"
	"fmt"
	"context"

	"github.com/hashicorp/packer/command"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/helper/multistep"
)

type StepBuilder struct {
	BuilderConfig map[string]interface{}
	builder       packer.Builder
}

func (s *StepBuilder) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	err := s.doRun(state)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	return multistep.ActionContinue
}

func (s *StepBuilder) doRun(state multistep.StateBag) error {
	ui := state.Get("ui").(packer.Ui)
	hook := state.Get("hook").(packer.Hook)
	cache := state.Get("cache").(packer.Cache)
	sourcePath := state.Get("source_path").(string)

	ui.Message("(vagrant) Builder type ...")
	builderType, err := s.builderType()
	if err != nil {
		return err
	}
	ui.Message(fmt.Sprintf("(vagrant) Builder type: %s", builderType))

	ui.Message("(vagrant) Builder ...")
	builder, found := command.Builders[builderType]
	if !found {
		return fmt.Errorf("unsupported builder type: %s", builderType)
	}
	ui.Message("(vagrant) Builder: OK")

	ui.Message("(vagrant) Builder prepare ...")
	s.builder = builder
	s.BuilderConfig["source_path"] = sourcePath
	warnings, err := s.builder.Prepare(s.BuilderConfig)
	if err != nil {
		return err
	}
	if warnings != nil && len(warnings) > 0 {
		for _, w := range warnings {
			ui.Message(fmt.Sprintf("(vagrant) WARNING: %s", w))
		}
	}
	ui.Message("(vagrant) Builder prepare: OK")

	ui.Message("(vagrant) Builder run ...")
	a, err := s.builder.Run(ui, hook, cache)
	if err != nil {
		return err
	}
	ui.Message("(vagrant) Builder run: OK")
	state.Put("artifact", a)
	return nil
}

func (s *StepBuilder) Cleanup(state multistep.StateBag) {
	s.builder.Cancel()
}

func (s *StepBuilder) builderType() (string, error) {
	raw, ok := s.BuilderConfig["type"]
	if !ok {
		return "", errors.New("invalid builder config, missing type")
	}

	t, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("invalid builder type value: %#v", t)
	}

	return t, nil
}

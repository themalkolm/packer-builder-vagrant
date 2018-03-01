package main

import (
	"errors"
	"log"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/helper/multistep"
)

type Builder struct {
	config *Config
	runner multistep.Runner
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

	//
	// Quick fix to pass -force flag down to the wrapped builder
	//
	// https://github.com/themalkolm/packer-builder-vagrant/issues/2
	//
	b.config.BuilderConfig["packer_force"] = b.config.PackerForce
	b.config.BuilderConfig["packer_debug"] = b.config.PackerDebug
	b.config.BuilderConfig["packer_on_error"] = b.config.PackerOnError
	b.config.BuilderConfig["packer_user_variables"] = b.config.PackerUserVars

	return warnings, nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	// Set up the state.
	state := new(multistep.BasicStateBag)
	state.Put("ui", ui)
	state.Put("hook", hook)
	state.Put("cache", cache)

	// Build the steps.
	steps := []multistep.Step{
		&StepFetchBox{
			URL:           b.config.URL,
			Name:          b.config.Name,
			Version:       b.config.Version,
			Provider:      b.config.Provider,
			BoxFile:       b.config.BoxFile,
			BuilderConfig: b.config.BuilderConfig,
		},
		&StepBuilder{
			BuilderConfig: b.config.BuilderConfig,
		},
	}

	// Run the steps.
	b.runner = common.NewRunnerWithPauseFn(steps, b.config.PackerConfig, ui, state)
	b.runner.Run(state)

	// Report any errors.
	if err, ok := state.GetOk("error"); ok {
		return nil, err.(error)
	}

	// If we were interrupted or cancelled, then just exit.
	if _, ok := state.GetOk(multistep.StateCancelled); ok {
		return nil, errors.New("Build was cancelled.")
	}

	if _, ok := state.GetOk(multistep.StateHalted); ok {
		return nil, errors.New("Build was halted.")
	}

	return state.Get("artifact").(packer.Artifact), nil
}

func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}

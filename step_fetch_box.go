package main

import (
	"fmt"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
)

type StepFetchBox struct {
	URL      string
	Name     string
	Version  string
	Provider string
	BoxFile  string

	BuilderConfig map[string]interface{}
}

func (s *StepFetchBox) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	err := s.doRun(state)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepFetchBox) doRun(state multistep.StateBag) error {
	ui := state.Get("ui").(packer.Ui)

	ui.Message("(vagrant) Builder source_path ...")
	if sourcePath, ok := s.BuilderConfig["source_path"]; ok {
		ui.Message(fmt.Sprintf("(vagrant) Builder source_path: %s", sourcePath))
		state.Put("source_path", sourcePath)
		return nil
	}

	v, err := NewVagrant(ui)
	if err != nil {
		return err
	}

	sourcePath, err := v.fetchBoxFile(s.URL, s.Name, s.Version, s.Provider, s.BoxFile)
	if err != nil {
		return err
	}

	ui.Message(fmt.Sprintf("(vagrant) Builder source_path: %s", sourcePath))
	state.Put("source_path", sourcePath)
	return nil
}

func (s *StepFetchBox) Cleanup(state multistep.StateBag) {
}

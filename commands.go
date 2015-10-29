package main

import (
	"os"

	"github.com/cskksc/sr6/command"
	"github.com/cskksc/sr6/command/agent"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available sr6 commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}
	Commands = map[string]cli.CommandFactory{
		"join": func() (cli.Command, error) {
			return &command.JoinCommand{
				Ui: ui,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Version: Version,
				Ui:      ui,
			}, nil
		},
		"agent": func() (cli.Command, error) {
			return &agent.Command{
				Ui: ui,
			}, nil
		},
	}
}

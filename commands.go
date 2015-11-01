package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cskksc/sr6/command"
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
			return &command.AgentCommand{
				Ui:         ui,
				ShutdownCh: makeShutdownCh(),
			}, nil
		},
		"members": func() (cli.Command, error) {
			return &command.MembersCommand{
				Ui: ui,
			}, nil
		},
		"leave": func() (cli.Command, error) {
			return &command.LeaveCommand{
				Ui: ui,
			}, nil
		},
	}
}

// makeShutdownCh returns a channel that can be used for shutdown
// notifications for commands. This channel will send a message for every
// interrupt or SIGTERM received.
func makeShutdownCh() <-chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()

	return resultCh
}

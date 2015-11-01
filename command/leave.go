package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/cskksc/sr6/sr6"
	"github.com/mitchellh/cli"
)

type LeaveCommand struct {
	Ui cli.Ui
}

func (c *LeaveCommand) Help() string {
	helpText := `
Usage: sr6 leave ...

  Causes the agent to gracefully leave the sr6 cluster and shutdown.
`
	return strings.TrimSpace(helpText)
}

func (c *LeaveCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("leave", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	client, err := sr6.NewRPCClient("localhost:8300")
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error connecting to sr6 agent: %s", err))
		return 1
	}
	defer client.Close()

	if err := client.Leave(); err != nil {
		c.Ui.Error(fmt.Sprintf("Error leaving cluster: %s", err))
		return 1
	}
	return 0
}

func (c *LeaveCommand) Synopsis() string {
	return "Gracefully leaves the sr6 cluster"
}

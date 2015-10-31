package command

import (
	"log"
	"strings"

	"github.com/cskksc/sr6/sr6"
	"github.com/mitchellh/cli"
)

type JoinCommand struct {
	Ui cli.Ui
}

func (c *JoinCommand) Help() string {
	helpText := `
Usage: sr6 join address ...
  Tells a running sr6 agent (with "sr6 agent") to join the cluster
  by specifying at least one existing member.
`
	return strings.TrimSpace(helpText)
}

func (c *JoinCommand) Run(args []string) int {
	addrs := args
	if len(addrs) == 0 {
		c.Ui.Error("At least one address to join must be specified.")
		c.Ui.Error("")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := sr6.NewRPCClient("localhost:8300")
	if err != nil {
		log.Println(err)
		return 1
	}
	defer client.Close()

	client.Join(addrs)
	return 0
}

func (c *JoinCommand) Synopsis() string {
	return "Tell sr6 agent to join cluster"
}

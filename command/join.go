package command

import (
	"log"
	"strings"

	"github.com/cskksc/sr6/command/agent"
	"github.com/mitchellh/cli"
)

type JoinCommand struct {
	Ui cli.Ui
}

func (c *JoinCommand) Help() string {
	helpText := `
Usage: sr6 join [options] address ...
  Tells a running sr6 agent (with "sr6 agent") to join the cluster
  by specifying at least one existing member.
Options:
  -rpc-addr=127.0.0.1:8400  RPC address of the sr6 agent.
`
	return strings.TrimSpace(helpText)
}

func (c *JoinCommand) Run(args []string) int {
	client, err := agent.NewRPCClient("localhost:8300")
	if err != nil {
		log.Println(err)
		return 1
	}
	client.Join()
	return 0
}

func (c *JoinCommand) Synopsis() string {
	return "Tell sr6 agent to join cluster"
}

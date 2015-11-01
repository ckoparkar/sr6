package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/cskksc/sr6/sr6"
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

  -rpc-addr=127.0.0.1:8300 RPC address of the sr6 agent.
`
	return strings.TrimSpace(helpText)
}

func (c *JoinCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("join", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	rpcAddr := RPCAddrFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	addrs := cmdFlags.Args()
	if len(addrs) == 0 {
		c.Ui.Error("At least one address to join must be specified.")
		c.Ui.Error("")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := sr6.NewRPCClient(*rpcAddr)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error connecting to sr6 agent: %s", err))
		return 1
	}
	defer client.Close()

	n, err := client.Join(addrs)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error joining the cluster: %s", err))
		return 1
	}
	c.Ui.Output(fmt.Sprintf("Successfully joined cluster by contacting %d nodes.", n))
	return 0
}

func (c *JoinCommand) Synopsis() string {
	return "Tell sr6 agent to join cluster"
}

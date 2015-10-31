package command

import (
	"flag"
	"log"
	"strings"

	"github.com/cskksc/sr6/sr6"
	"github.com/mitchellh/cli"
)

type MembersCommand struct {
	Ui cli.Ui
}

func (c *MembersCommand) Help() string {
	helpText := `Usage: sr6 members

  Outputs the members connected to the running sr6 agent.`
	return strings.TrimSpace(helpText)
}

func (c *MembersCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("agent", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	client, err := sr6.NewRPCClient("localhost:8300")
	if err != nil {
		log.Println(err)
		return 1
	}
	defer client.Close()

	client.Members()
	return 0
}

func (c *MembersCommand) Synopsis() string {
	return "Lists the members of a sr6 cluster"
}

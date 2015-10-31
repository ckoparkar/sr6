package command

import (
	"log"

	"github.com/cskksc/sr6/command/agent"
	"github.com/mitchellh/cli"
)

type MembersCommand struct {
	Ui cli.Ui
}

func (c *MembersCommand) Help() string {
	return ""
}

func (c *MembersCommand) Run(args []string) int {
	client, err := agent.NewRPCClient("localhost:8300")
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

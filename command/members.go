package command

import (
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/cskksc/sr6/sr6"
	"github.com/hashicorp/serf/serf"
	"github.com/mitchellh/cli"
	"github.com/ryanuber/columnize"
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
		c.Ui.Error(fmt.Sprintf("Error connecting to sr6 agent: %s", err))
		return 1
	}
	defer client.Close()

	members, err := client.Members()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error getting cluster members: %s", err))
		return 1
	}
	result := columnize.SimpleFormat(output(members))
	fmt.Println(result)
	return 0
}

func (c *MembersCommand) Synopsis() string {
	return "Lists the members of a sr6 cluster"
}

func output(members []serf.Member) []string {
	result := make([]string, 0, len(members))
	header := "Node|Address|Status|Type"
	result = append(result, header)
	for _, member := range members {
		addr := net.TCPAddr{IP: member.Addr, Port: int(member.Port)}
		line := fmt.Sprintf("%s|%s|%s|%s",
			member.Name, addr.String(), member.Status, member.Tags["role"])
		result = append(result, line)
	}
	return result
}

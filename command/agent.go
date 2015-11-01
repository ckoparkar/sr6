package command

import (
	"flag"
	"log"
	"strings"

	"github.com/cskksc/sr6/sr6"
	"github.com/mitchellh/cli"
)

type AgentCommand struct {
	Ui         cli.Ui
	ShutdownCh <-chan struct{}
	args       []string
	server     *sr6.Server
}

func (c *AgentCommand) Help() string {
	helpText := `Usage: sr6 agent [options]

  Starts the sr6 agent and runs until an interrupt is received. The
  agent represents a single node in a cluster.`
	return strings.TrimSpace(helpText)
}

func (c *AgentCommand) Run(args []string) int {
	c.args = args
	config, err := c.readConfig()
	if err != nil {
		log.Fatal(err)
	}
	s, err := sr6.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	c.server = s
	return c.handleSignals()
}

func (c *AgentCommand) Synopsis() string {
	return "Runs a sr6 agent"
}

// Runs in its own go routine
// handleSignals monitors the shutdownCh channel and acts on it
func (c *AgentCommand) handleSignals() int {
	select {
	case <-c.ShutdownCh:
		if err := c.server.Shutdown(); err != nil {
			log.Println("[INFO] sr6: Couldn't properly shutdown the server")
			return 1
		}
		return 0
	}
}

// readConfig reads config provided as cmd-line args,
// and merges it with the defaults
func (c *AgentCommand) readConfig() (*sr6.Config, error) {
	var cmdConfig sr6.Config
	cmdFlags := flag.NewFlagSet("agent", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&cmdConfig.NodeName, "node", "", "node name")
	cmdFlags.BoolVar(&cmdConfig.Leader, "leader", false, "enable server leader node")
	if err := cmdFlags.Parse(c.args); err != nil {
		return nil, err
	}

	config, err := sr6.DefaultConfig()
	if err != nil {
		return nil, err
	}
	// Not all config would be provided as cmd-line args
	config = sr6.MergeConfig(config, &cmdConfig)
	return config, nil
}

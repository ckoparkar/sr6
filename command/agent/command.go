package agent

import (
	"flag"
	"log"

	"github.com/mitchellh/cli"
)

type Command struct {
	Ui         cli.Ui
	ShutdownCh <-chan struct{}
	args       []string
	server     *Server
}

func (c *Command) Help() string {
	return ""
}

func (c *Command) Run(args []string) int {
	c.args = args
	config, err := c.readConfig()
	if err != nil {
		log.Fatal(err)
	}
	s, err := NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	c.server = s
	return c.handleSignals()
}

func (c *Command) Synopsis() string {
	return "Start a sr6 agent"
}

// Runs in its own go routine
// handleSignals monitors the shutdownCh channel and acts on it
func (c *Command) handleSignals() int {
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
func (c *Command) readConfig() (*Config, error) {
	var cmdConfig Config
	cmdFlags := flag.NewFlagSet("agent", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&cmdConfig.NodeName, "node", "", "node name")
	if err := cmdFlags.Parse(c.args); err != nil {
		return nil, err
	}

	config, err := DefaultConfig()
	if err != nil {
		return nil, err
	}
	// Not all config would be provided as cmd-line args
	config = MergeConfig(config, &cmdConfig)
	return config, nil
}

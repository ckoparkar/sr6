package agent

import (
	"flag"
	"log"
	"net/rpc"

	"github.com/mitchellh/cli"
)

type Command struct {
	Ui   cli.Ui
	args []string
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
	s.rpcServer.Accept(s.rpcListener)
	for {
		conn, err := s.rpcListener.Accept()
		if err != nil {
			log.Println(err)
		}
		rpc.ServeConn(conn)
	}
	return 0
}

func (c *Command) Synopsis() string {
	return "Start a sr6 agent"
}

func (c *Command) readConfig() (*Config, error) {
	cmdFlags := flag.NewFlagSet("agent", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	config, err := DefaultConfig()
	if err != nil {
		return nil, err
	}
	cmdFlags.StringVar(&config.NodeName, "node", "", "node name")
	if err := cmdFlags.Parse(c.args); err != nil {
		return nil, err
	}

	config.SerfConfig.NodeName = config.NodeName
	return config, nil
}

package command

import (
	"fmt"

	"github.com/mitchellh/cli"
)

type VersionCommand struct {
	Version string
	Ui      cli.Ui
}

func (c *VersionCommand) Help() string {
	return ""
}

func (c *VersionCommand) Run(_ []string) int {
	fmt.Printf("sr6: %s\n", c.Version)
	return 0
}

func (c *VersionCommand) Synopsis() string {
	return "Prints the sr6 version"
}

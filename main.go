package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

// var (
//	nodeName = flag.String("nodename", "hostname", "Advertise this nodename.")
//	mode     = flag.String("mode", "server", "Decides whether to run as client/server.")
//	listen   = flag.String("listen", ":8080", "HTTP listen address.")
// )

func main() {
	args := os.Args[1:]
	cli := &cli.CLI{
		Args:     args,
		Commands: Commands,
		HelpFunc: cli.BasicHelpFunc("sr6"),
	}
	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(exitCode)
}

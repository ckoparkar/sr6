package agent

import "fmt"

type endpoints struct {
	Internal *Internal
}

type Internal struct {
	srv *Server
}

func (i *Internal) Join(args string, reply *string) error {
	fmt.Println("I am here", args)
	*reply = "I will join this node."
	return nil
}

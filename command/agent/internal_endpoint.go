package agent

import "fmt"

type endpoints struct {
	Internal *Internal
}

// Internal serves as a endpoint for all internal operations
// These API's may not be directly exposed to clients
type Internal struct {
	srv *Server
}

func (i *Internal) Join(args string, reply *string) error {
	fmt.Println("I am here", args)
	*reply = "I will join this node."
	return nil
}

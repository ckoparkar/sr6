package agent

import "github.com/hashicorp/serf/serf"

type endpoints struct {
	Internal *Internal
}

// Internal serves as a endpoint for all internal operations
// These API's may not be directly exposed to clients
type Internal struct {
	srv *Server
}

func (i *Internal) Join(addrs []string, reply *int) error {
	n, err := i.srv.serfLAN.Join(addrs, true)
	if err != nil {
		return err
	}
	*reply = n
	return nil
}

func (i *Internal) Members(_ string, reply *[]serf.Member) error {
	*reply = i.srv.serfLAN.Members()
	return nil
}

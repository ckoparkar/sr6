package sr6

import (
	"net/rpc"

	"github.com/hashicorp/serf/serf"
)

type RPCClient struct {
	conn *rpc.Client
}

func NewRPCClient(addr string) (*RPCClient, error) {
	conn, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &RPCClient{
		conn: conn,
	}, nil
}

func (c *RPCClient) Join(addrs []string) (int, error) {
	var reply int
	if err := c.conn.Call("Internal.Join", addrs, &reply); err != nil {
		return -1, err
	}
	return reply, nil
}

func (c *RPCClient) Members() ([]serf.Member, error) {
	var reply []serf.Member
	if err := c.conn.Call("Internal.Members", "", &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *RPCClient) Close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}

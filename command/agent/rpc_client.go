package agent

import (
	"fmt"
	"log"
	"net/rpc"
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

func (c *RPCClient) Join(addrs []string) {
	var reply int
	if err := c.conn.Call("Internal.Join", addrs, &reply); err != nil {
		log.Println(err)
	}
	fmt.Println(reply)
}

func (c *RPCClient) Close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}

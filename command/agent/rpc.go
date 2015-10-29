package agent

import (
	"fmt"
	"log"
	"net/rpc"
)

type joinRequest struct {
	Existing []string
}

type joinResponse struct {
	Num int32
}

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

func (c *RPCClient) Join() {
	var reply string
	if err := c.conn.Call("Internal.Join", "123", &reply); err != nil {
		log.Println(err)
	}
	fmt.Println(reply)
}

package agent

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"github.com/hashicorp/serf/serf"
)

const (
	DefaultSerfPort = 8201
)

var (
	DefaultRPCAddr = &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8300}
)

type Config struct {
	// Node name is the name we use to advertise. Defaults to hostname.
	NodeName string

	// SerfConfig is the configuration for serf
	SerfConfig *serf.Config

	// RPCAddr describes the
	RPCAddr *net.TCPAddr
}

func DefaultConfig() (*Config, error) {
	eventCh := make(chan serf.Event, 256)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	c := &Config{
		NodeName:   hostname,
		SerfConfig: serf.DefaultConfig(),
		RPCAddr:    DefaultRPCAddr,
	}

	// Serf config
	c.SerfConfig.NodeName = hostname
	c.SerfConfig.EventCh = eventCh
	c.SerfConfig.SnapshotPath = "/tmp/serf.snapshot"
	c.SerfConfig.MemberlistConfig.BindPort = DefaultSerfPort
	c.SerfConfig.MemberlistConfig.AdvertisePort = DefaultSerfPort

	return c, nil
}

type Server struct {
	config  *Config
	serfLAN *serf.Serf

	// rpcListener is used to listen for incoming connections
	rpcListener net.Listener
	rpcServer   *rpc.Server

	// endpoints holds our RPC endpoints
	endpoints endpoints
}

func NewServer(config *Config) (*Server, error) {
	serfLAN, err := serf.Create(config.SerfConfig)
	if err != nil {
		return nil, err
	}
	s := &Server{
		config:    config,
		serfLAN:   serfLAN,
		rpcServer: rpc.NewServer(),
	}
	if err := s.setupRPC(); err != nil {
		log.Fatal(err)
	}

	return s, nil
}

func (s *Server) setupRPC() error {
	s.endpoints.Internal = &Internal{s}
	if err := s.rpcServer.Register(s.endpoints.Internal); err != nil {
		return err
	}

	list, err := net.ListenTCP("tcp", s.config.RPCAddr)
	if err != nil {
		return err
	}
	s.rpcListener = list
	return nil
}

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

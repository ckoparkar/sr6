package agent

import (
	"log"
	"net"
	"net/rpc"

	"github.com/hashicorp/serf/serf"
)

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

func (s *Server) Shutdown() error {
	log.Printf("[INFO] sr6: shutting down server")
	if s.serfLAN != nil {
		s.serfLAN.Leave()
		s.serfLAN.Shutdown()
	}
	if s.rpcListener != nil {
		s.rpcListener.Close()
	}
	return nil
}

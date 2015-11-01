package sr6

import (
	"log"
	"net"
	"net/rpc"
	"sync"

	"github.com/hashicorp/serf/serf"
)

// Server is the main sr6 server
type Server struct {
	// config is this servers config
	config *Config

	// serfLAN is the local serf cluster
	serfLAN *serf.Serf

	// eventChLAN is used to receive events from the
	// serf cluster in the datacenter
	eventChLAN chan serf.Event

	// rpcListener is used to listen for incoming connections
	rpcListener net.Listener
	rpcServer   *rpc.Server

	// endpoints holds our RPC endpoints
	endpoints endpoints

	// clean studown
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

func NewServer(config *Config) (*Server, error) {
	s := &Server{
		config:     config,
		eventChLAN: make(chan serf.Event, 256),
		rpcServer:  rpc.NewServer(),
		shutdownCh: make(chan struct{}),
	}

	// Setup serf
	serfLAN, err := s.setupSerf()
	if err != nil {
		return nil, err
	}
	s.serfLAN = serfLAN

	// Setup RPC and start listening for requests
	if err := s.setupRPC(); err != nil {
		log.Fatal(err)
	}
	go s.listenRPC()

	return s, nil
}

// setupSerf sets up serf and provides a handle on its events
func (s *Server) setupSerf() (*serf.Serf, error) {
	// initializes the config (contains maps)
	s.config.SerfConfig.Init()
	s.config.SerfConfig.EventCh = s.eventChLAN
	if s.config.Leader {
		s.config.SerfConfig.Tags["role"] = "leader"
	} else {
		s.config.SerfConfig.Tags["role"] = "follower"
	}
	serfLAN, err := serf.Create(s.config.SerfConfig)
	if err != nil {
		return nil, err
	}
	return serfLAN, nil
}

// setupRPC starts a RPC server and registers all endpoints
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

// listenRPC serves all incoming RPC requests
func (s *Server) listenRPC() {
	s.rpcServer.Accept(s.rpcListener)
	for {
		conn, err := s.rpcListener.Accept()
		if err != nil {
			if s.shutdown {
				return
			}
			log.Printf("[ERR] sr6.rpc: failed to accept RPC conn: %v", err)
		}
		log.Printf("[INFO] rpc: Accepted client: %v", conn.RemoteAddr())
		rpc.ServeConn(conn)
	}
}

// Shutdown closes all active servers running in background
// this method is called when Ctrl+C signal is received on shutdownCh
func (s *Server) Shutdown() error {
	log.Printf("[INFO] sr6: shutting down server")
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		return nil
	}
	s.shutdown = true

	if s.serfLAN != nil {
		s.serfLAN.Leave()
		s.serfLAN.Shutdown()
	}
	if s.rpcListener != nil {
		s.rpcListener.Close()
	}
	return nil
}

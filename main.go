package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/serf/serf"
)

var (
	nodeName = flag.String("nodename", "hostname", "Advertise this nodename.")
	join     = flag.String("join", "localhost:8202", "Attempts to join `node`.")
	port     = flag.Int("port", 8201, "Port used for serf bind/advertise.")
)

func main() {
	flag.Parse()
	c := DefaultConfig()
	s, err := NewServer(c)
	if err != nil {
		log.Fatal(err)
	}
	for {
		_, err := s.serfLAN.Join([]string{*join}, true)
		if err != nil {
			log.Println(err)
		} else {
			break
		}
		time.Sleep(2 * time.Second)
	}
	for {
		fmt.Println(s.serfLAN.Members())
		time.Sleep(5 * time.Second)
	}
}

type Config struct {
	// Node name is the name we use to advertise. Defaults to hostname.
	NodeName string

	// SerfConfig is the configuration for serf
	SerfConfig *serf.Config
}

func DefaultConfig() *Config {
	eventCh := make(chan serf.Event, 256)
	c := &Config{
		NodeName:   *nodeName,
		SerfConfig: serf.DefaultConfig(),
	}
	c.SerfConfig.NodeName = *nodeName
	c.SerfConfig.EventCh = eventCh
	c.SerfConfig.SnapshotPath = "/tmp/serf.snapshot"
	c.SerfConfig.MemberlistConfig.BindPort = *port
	c.SerfConfig.MemberlistConfig.AdvertisePort = *port

	return c
}

type Server struct {
	serfLAN *serf.Serf
}

func NewServer(config *Config) (*Server, error) {
	serfLAN, err := serf.Create(config.SerfConfig)
	if err != nil {
		return nil, err
	}
	return &Server{
		serfLAN: serfLAN,
	}, nil
}

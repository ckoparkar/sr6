package main

import (
	"github.com/hashicorp/serf/serf"
)

var (
	serfPort = 8201

	nodeName = "abc"
)

type Config struct {
	// Node name is the name we use to advertise. Defaults to hostname.
	NodeName string

	// SerfConfig is the configuration for serf
	SerfConfig *serf.Config
}

func DefaultConfig() *Config {
	eventCh := make(chan serf.Event, 256)
	c := &Config{
		NodeName:   nodeName,
		SerfConfig: serf.DefaultConfig(),
	}
	c.SerfConfig.NodeName = nodeName
	c.SerfConfig.EventCh = eventCh
	c.SerfConfig.SnapshotPath = "/tmp/serf.snapshot"
	c.SerfConfig.MemberlistConfig.BindPort = serfPort
	c.SerfConfig.MemberlistConfig.AdvertisePort = serfPort

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

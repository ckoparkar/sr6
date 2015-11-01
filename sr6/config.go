package sr6

import (
	"net"
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

	// Leader decides if we are the leader
	Leader bool

	// SerfConfig is the configuration for serf
	SerfConfig *serf.Config

	// RPCAddr describes the
	RPCAddr *net.TCPAddr
}

func DefaultConfig() (*Config, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	c := &Config{
		NodeName:   hostname,
		Leader:     false,
		SerfConfig: serf.DefaultConfig(),
		RPCAddr:    DefaultRPCAddr,
	}
	c.setupSerfConfig()
	return c, nil
}

func (c *Config) setupSerfConfig() error {
	eventCh := make(chan serf.Event, 256)
	c.SerfConfig.NodeName = c.NodeName
	c.SerfConfig.EventCh = eventCh
	c.SerfConfig.SnapshotPath = "/tmp/serf.snapshot"
	c.SerfConfig.MemberlistConfig.BindPort = DefaultSerfPort
	c.SerfConfig.MemberlistConfig.AdvertisePort = DefaultSerfPort
	return nil
}

func MergeConfig(a, b *Config) *Config {
	var result Config = *a
	if b.NodeName != "" {
		result.NodeName = b.NodeName
		result.SerfConfig.NodeName = b.NodeName
	}
	if b.Leader {
		result.Leader = b.Leader
	}
	return &result
}

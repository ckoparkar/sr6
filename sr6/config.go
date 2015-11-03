package sr6

import (
	"net"
	"os"
	"time"

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

	// HostsFile points to hosts file. Defaults to /etc/hosts
	HostsFile string

	// HostSuffix is os.Hostname suffix
	HostSuffix string

	// HostsUpdateInterval decides when to update hosts file
	HostsUpdateInterval time.Duration
}

func DefaultConfig() (*Config, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	c := &Config{
		NodeName:   hostname,
		SerfConfig: serf.DefaultConfig(),
		RPCAddr:    DefaultRPCAddr,
		HostsFile:  "/etc/hosts",

		HostsUpdateInterval: 10 * time.Second,
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
	if b.HostsFile != "" {
		result.HostsFile = b.HostsFile
	}
	if b.HostSuffix != "" {
		result.HostSuffix = b.HostSuffix
	}
	return &result
}

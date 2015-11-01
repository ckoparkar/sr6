package sr6

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type HostsManager struct {
	// it is the map of ip -> hostname
	hosts map[string]string

	// path of hosts file. defaults to /etc/hosts
	path string

	sync.Mutex
}

// NewHosts parses hosts file at *path*
func NewHostsManager(path string) (*HostsManager, error) {
	h := &HostsManager{
		hosts: make(map[string]string),
		path:  path,
	}
	// if hosts file doesnt exist, return
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("[WARN] writing a new hosts file at: %s", path)
		return h, nil
	}
	// parse hosts file
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(input), "\n")
	for _, l := range lines {
		x := strings.Split(l, " ")
		if len(x) < 2 {
			continue
		}
		ip := strings.TrimSpace(x[0])
		hostname := strings.TrimSpace(x[1])
		if len(x) == 3 {
			hostname = x[2]
		}
		h.hosts[ip] = hostname
	}
	return h, nil
}

// add the ip-hostname pair and write changes in file
func (h *HostsManager) add(ip, hostname string) error {
	h.Lock()
	defer h.Unlock()
	h.hosts[ip] = hostname
	var content string
	for k, v := range h.hosts {
		content += fmt.Sprintf("%s %s\n", k, v)
	}
	f, err := os.OpenFile(h.path, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(content)
	return nil
}

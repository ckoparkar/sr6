package sr6

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// NewHosts parses hosts file at *path*
func NewHosts(path string) (map[string]string, error) {
	// re := regexp.MustCompile("(.*) +(.*)")
	hosts := make(map[string]string)
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
		hosts[ip] = hostname
	}
	return hosts, nil
}

func (s *Server) updateHosts(ip, hostname string) error {
	s.hostsLock.Lock()
	defer s.hostsLock.Unlock()

	s.hosts[ip] = hostname
	for k, v := range s.hosts {
		fmt.Println(k, "---", v)
	}
	// TODO(cskksc): update hosts file here
	return nil
}

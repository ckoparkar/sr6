package types

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"time"

	sigar "github.com/cloudfoundry/gosigar"
	"github.com/cskksc/sr6/request"
)

var ErrNotFound = fmt.Errorf("Not found.")

type SSHKeys struct {
	Private, Public string
}

func NewSSHKeys(sshPrivateKeyPath, sshPublicKeyPath string) (*SSHKeys, error) {
	privateKey, err := ioutil.ReadFile(sshPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := ioutil.ReadFile(sshPublicKeyPath)
	if err != nil {
		return nil, err
	}
	return &SSHKeys{
		Private: string(privateKey),
		Public:  string(publicKey),
	}, nil
}

type BaseResponse struct {
	Payload interface{} `json:"payload"`

	Status  int    `json:"status"`
	Message string `json:"message"`
}

type RegisterResponse struct {
	SSHKeys      SSHKeys `json:"ssh_keys"`
	Hostname     string  `json:"hostname"`
	PollInterval string  `json:"poll_interval"`
	Status       int     `json:"status"`
	Message      string  `json:"message"`
}

type HeartbeatResponse struct {
	*Follower `json:"follower"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
}

type Follower struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	MemUsed string `json:"mem_used"`
}

func NewFollower(id, listenAddr string) (*Follower, error) {
	ip, err := internalIP()
	if err != nil {
		return nil, err
	}
	memUsed := memUsage()
	return &Follower{
		ID:      id,
		Address: ip + listenAddr,
		MemUsed: memUsed,
	}, nil
}

func (f *Follower) Ping() error {
	ping := make(chan int)
	go func() {
		req := request.NewRequest("GET", "http", f.Address, "/heartbeat", nil, nil)
		_, _, err := req.Do()
		if err == nil {
			ping <- 1
		}
	}()

	select {
	case <-ping:
	case <-time.After(10 * time.Millisecond):
		return ErrNotFound
	}
	return nil
}

func internalIP() (string, error) {
	re := regexp.MustCompile("[0-9]+.[0-9]+.[0-9]+.[0-9]+")
	name, err := os.Hostname()
	if err != nil {
		log.Printf("Couldn't get IP, %v", err)
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		log.Printf("Couldn't get IP, %v", err)
	}
	for _, a := range addrs {
		if ip := re.FindString(a); ip != "" {
			return ip, nil
		}
	}

	return "", ErrNotFound
}

func memUsage() string {
	mem := sigar.Mem{}
	mem.Get()
	used := float64(mem.ActualUsed) / (float64(mem.ActualFree) + float64(mem.ActualUsed)) * 100
	return fmt.Sprintf("%.2f", used)
}

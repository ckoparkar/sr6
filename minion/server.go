package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"sync"
	"time"

	"github.com/cskksc/sr6/request"
	"github.com/cskksc/sr6/types"
)

type Server struct {
	ID string

	restartMonitor chan int
	pollInterval   time.Duration

	mu       sync.RWMutex
	lastPoll time.Time
}

func NewServer() *Server {
	id, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	s := &Server{
		ID: id,

		restartMonitor: make(chan int),
		pollInterval:   time.Minute,
	}
	go s.monitor()
	return s
}

// Runs in its own go routine
// If last poll request from master was over *pollInterval* ago,
// try to re-register (depends on poll interval of master)
func (s *Server) monitor() {
	ticker := time.NewTicker(s.pollInterval)
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			// if we dint receive poll for 5 cycles, re-register
			if now.Sub(s.lastPoll) > (s.pollInterval * 2) {
				if err := s.register(); err != nil {
					// If we cannot re-register, bail out
					log.Fatal(err)
				}
			}
		case <-s.restartMonitor:
			return
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/heartbeat" {
		s.serveHeartbeat(w, r)
		return
	}
}

func (s *Server) serveHeartbeat(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastPoll = time.Now()

	f, err := types.NewFollower(s.ID, *listenAddr)
	if err != nil {
		s.serveError(w, r, err, http.StatusNotFound)
		return
	}
	resp := types.HeartbeatResponse{
		Follower: f,
		Status:   http.StatusOK,
		Message:  "success",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) serveError(w http.ResponseWriter, r *http.Request, err error, status int) {
	resp := types.BaseResponse{
		Status:  status,
		Message: err.Error(),
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) register() error {
	f, err := types.NewFollower(s.ID, *listenAddr)
	if err != nil {
		return err
	}
	buf, err := json.Marshal(f)
	if err != nil {
		return err
	}
	req := request.NewRequest("POST", "http", *masterAddr, "/register", bytes.NewReader(buf), nil)
	_, resp, err := req.Do()
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var rr types.RegisterResponse
	json.Unmarshal(body, &rr)

	// Change our hostname, add ssh keys to proper files
	// TODO(cskksc): setup a docker test env
	// if err := setHostname(rr.Hostname); err != nil {
	//	return err
	// }

	if err := writeSSHKeys(rr.SSHKeys); err != nil {
		return err
	}

	// Set our poll interval
	pi, err := time.ParseDuration(rr.PollInterval)
	if s.pollInterval != pi && err == nil {
		s.pollInterval = pi
		s.restartMonitor <- 1
		go s.monitor()
	}

	log.Printf("Registered successfully with %s.\n", *masterAddr)
	return nil
}

func writeSSHKeys(keys types.SSHKeys) error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	publicKeyPath := currentUser.HomeDir + "/.ssh/id_rsa.pub"
	privateKeyPath := currentUser.HomeDir + "/.ssh/id_rsa"
	if err := ioutil.WriteFile(publicKeyPath, []byte(keys.Public), 0644); err != nil {
		return err
	}
	if err := ioutil.WriteFile(privateKeyPath, []byte(keys.Private), 0644); err != nil {
		return err
	}
	return nil
}

func setHostname(name string) error {
	cmd := exec.Command("sudo", "hostname", name)
	var out bytes.Buffer
	cmd.Stdout = &out

	var in bytes.Buffer
	cmd.Stdin = &in

	in.Write([]byte("bazzinga"))
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("in all caps: %q\n", out.String())
	return nil
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/cskksc/sr6/request"
	"github.com/cskksc/sr6/types"
)

type Server struct {
	ID string

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
	}
	go s.monitor()
	return s
}

// Runs in its own go routine
// If last poll request from master was over *pollInterval* ago,
// try to re-register (depends on poll interval of master)
func (s *Server) monitor() {
	ticker := time.NewTicker(*pollInterval)
	for range ticker.C {
		now := time.Now()
		// if we dint receive poll for 5 cycles, re-register
		if now.Sub(s.lastPoll) > (*pollInterval * 5) {
			if err := register(s.ID); err != nil {
				// If we cannot re-register, bail out
				log.Fatal(err)
			}
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

func register(id string) error {
	f, err := types.NewFollower(id, *listenAddr)
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
	fmt.Println(rr)

	// Change our hostname, add ssh keys to proper files
	// TODO(cskksc)

	log.Printf("Registered successfully with %s.\n", *masterAddr)
	return nil
}

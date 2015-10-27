package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cskksc/sr6/request"
	"github.com/cskksc/sr6/types"
)

type Server struct {
	sshKeys   types.SSHKeys
	mu        sync.RWMutex
	followers []types.Follower
}

func NewServer() (*Server, error) {
	keys, err := types.NewSSHKeys(*sshPrivateKeyPath, *sshPublicKeyPath)
	if err != nil {
		return nil, err
	}
	s := &Server{
		sshKeys:   *keys,
		followers: make([]types.Follower, 0),
	}
	go s.run()
	return s, nil
}

func (s *Server) run() {
	ticker := time.NewTicker(*pollInterval)
	for range ticker.C {
		s.poll()
		s.inspect()
	}
}

func (s *Server) inspect() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fmt.Println(s.followers)
}

func (s *Server) poll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	ping := make(chan int)
	for i := len(s.followers) - 1; i >= 0; i-- {
		f := s.followers[i]
		hostport := f.Address + *followerPort
		timeout := time.After(10 * time.Millisecond)
		go func() {
			req := request.NewRequest("GET", "http", hostport, "/heartbeat", nil, nil)
			_, _, err := req.Do()
			if err != nil {
				<-timeout
			} else {
				ping <- 1
			}
		}()

		select {
		case <-ping:
		case <-time.After(10 * time.Millisecond):
			// delete this element
			s.followers = append(s.followers[:i], s.followers[i+1:]...)
			log.Printf("De-registered %#v\n", f)
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/register" {
		s.serveRegisterFollower(w, r)
		return
	}
}

func (s *Server) serveRegisterFollower(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var f types.Follower
	json.Unmarshal(body, &f)

	re := regexp.MustCompile("[[:digit:]]+")
	lastID := -1
	// Check if follower already exists.
	s.mu.RLock()
	// cf -> currentFollower
	for _, cf := range s.followers {
		if cf.ID == f.ID {
			// capture ids here to determine next
			lastID, _ = strconv.Atoi(re.FindString(cf.ID))

			// follower is already present
			return
		}
	}
	s.mu.RUnlock()

	// We have a new follower. Send ssh and host information
	lastID++
	hostname := strings.Replace(*hostnamePattern, "ID", strconv.Itoa(lastID), -1)
	resp := types.RegisterResponse{
		Hostname:     hostname,
		SSHKeys:      s.sshKeys,
		PollInterval: pollInterval.String(),
		Status:       http.StatusOK,
	}
	json.NewEncoder(w).Encode(resp)

	// Add the follower to list
	s.mu.Lock()
	defer s.mu.Unlock()
	s.followers = append(s.followers, f)
	log.Printf("Registered %#v\n", f)
}

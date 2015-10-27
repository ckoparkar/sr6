package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/cskksc/minion/request"
)

const (
	followerPort = ":8282"
)

type Server struct {
	mu        sync.RWMutex
	followers []Heartbeat
}

func NewServer() *Server {
	s := &Server{
		followers: make([]Heartbeat, 0),
	}
	go s.run()
	return s
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
		hostport := f.Address + ":8080"
		go func() {
			req := request.NewRequest("GET", "http", hostport, "/heartbeat", nil, nil)
			_, _, err := req.Do()
			if err != nil {
				log.Println(err)
			}
			ping <- 1
		}()

		select {
		case <-ping:
		case <-time.After(10 * time.Millisecond):
			// delete this element
			s.followers = append(s.followers[:i], s.followers[i+1:]...)
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/register" {
		s.registerFollower(w, r)
		return
	}
}

func (s *Server) registerFollower(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var beat Heartbeat
	json.Unmarshal(body, &beat)

	s.mu.RLock()
	for _, f := range s.followers {
		if f.ID == beat.ID {
			// follower is already present
			return
		}
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.followers = append(s.followers, beat)
}

type Heartbeat struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	MemUsed string `json:"mem_used"`
}

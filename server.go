package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

type Server struct {
	mu        sync.RWMutex
	followers []Heartbeat
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

	mu.RLock()
	for f := range s.followers {
		if f.ID == beat.ID {
			// follower is already present
			return
		}
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	s.followers = append(s.followers, beat)
	return
}

type Heartbeat struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	MemUsed string `json:"mem_used"`
}

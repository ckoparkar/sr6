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

type Server struct {
	mu        sync.RWMutex
	followers []Heartbeat
}

func NewServer() *Server {
	s := &Server{
		followers: make([]Heartbeat, 0),
	}
	go s.poll()
	go s.inspect()
	return s
}

func (s *Server) inspect() {
	for {
		time.Sleep(10 * time.Second)
		s.mu.RLock()
		fmt.Println(s.followers)
		s.mu.RUnlock()
	}
}

func (s *Server) poll() {
	for {
		time.Sleep(10 * time.Second)
		s.mu.Lock()
		ping := make(chan int)
		for i := len(s.followers) - 1; i >= 0; i-- {
			f := s.followers[i]
			hostport := fmt.Sprintf("%s:8080", f.Address)
			go func() {
				req := request.NewRequest("GET", "http", hostport, "/heartbeat", nil, nil)
				rtt, resp, err := req.Do()
				if err != nil {
					log.Println(err)
				}
				fmt.Println(rtt)
				fmt.Println(resp)
				ping <- 1
			}()

			select {
			case <-ping:
			case <-time.After(10 * time.Millisecond):
				// delete this element
				s.followers = append(s.followers[:i], s.followers[i+1:]...)
			}
		}
		s.mu.Unlock()
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
	fmt.Println(beat)

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
	fmt.Println(s.followers)
}

type Heartbeat struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	MemUsed string `json:"mem_used"`
}

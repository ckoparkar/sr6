package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"

	sigar "github.com/cloudfoundry/gosigar"
)

var (
	ErrNotFound = fmt.Errorf("Not found.")
)

type Response struct {
	Payload interface{} `json:"payload"`

	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Server struct {
	ID string
}

func NewServer() *Server {
	id, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return &Server{
		ID: id,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/heartbeat" {
		s.serveHeartbeat(w, r)
		return
	}
}

func (s *Server) serveHeartbeat(w http.ResponseWriter, r *http.Request) {
	beat, err := NewHeartbeat(s.ID)
	if err != nil {
		s.serveError(w, r, err, http.StatusNotFound)
		return
	}
	resp := Response{
		Payload: beat,
		Status:  http.StatusOK,
		Message: "success",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) serveError(w http.ResponseWriter, r *http.Request, err error, status int) {
	resp := Response{
		Status:  status,
		Message: err.Error(),
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

type Heartbeat struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	MemUsed string `json:"mem_used"`
}

func NewHeartbeat(id string) (*Heartbeat, error) {
	ip, err := internalIP()
	if err != nil {
		return nil, err
	}
	memUsed := memUsage()
	return &Heartbeat{
		ID:      id,
		Address: ip,
		MemUsed: memUsed,
	}, nil
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
	return fmt.Sprintf("%.2f%%", used)
}

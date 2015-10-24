package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"regexp"
	"time"
)

var (
	ErrNotFound = fmt.Errorf("Not found.")
)

type Response struct {
	Payload interface{} `json:"payload"`

	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Heartbeat struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
}

type Server struct {
	ID int
}

func NewServer() *Server {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int() % 99
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
	ip, err := internalIP()
	if err != nil {
		s.serveError(w, r, err)
		return
	}
	beat := Heartbeat{
		ID:      s.ID,
		Address: ip,
	}
	resp := Response{
		Payload: beat,
		Status:  http.StatusOK,
		Message: "success",
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) serveError(w http.ResponseWriter, r *http.Request, err error) {
	resp := Response{
		Status:  http.StatusInternalServerError,
		Message: err.Error(),
	}
	json.NewEncoder(w).Encode(resp)
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

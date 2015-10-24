package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	s := new(Server)
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Server struct {
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	listenAddr = flag.String("listen", ":8282", "HTTP listen adddress.")
	masterAddr = flag.String("master-addr", "localhost:8281", "IP address of master.")
)

func main() {
	flag.Parse()
	s := NewServer()
	if err := s.register(); err != nil {
		// if we cannot register at start,
		// we cannot proceed
		log.Fatal(err)
	}
	http.Handle("/", s)
	fmt.Printf("Listening for requests on %s...\n", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

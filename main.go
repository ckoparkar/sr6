package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	listenAddr   = flag.String("listen", ":8282", "HTTP listen adddress.")
	pollInterval = flag.Duration("poll", 5*time.Minute, "Registers itself with master, every `t`.")
	masterAddr   = flag.String("master", "localhost", "IP address of master.")
)

func main() {
	flag.Parse()
	s := NewServer()
	if err := register(s.ID); err != nil {
		// if we cannot register at start,
		// we cannot proceed
		log.Fatal(err)
	}
	http.Handle("/", s)
	fmt.Printf("Listening for requests on %s...\n", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

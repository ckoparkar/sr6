package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	listenAddr = flag.String("listen", ":8080", "HTTP listen adddress.")
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

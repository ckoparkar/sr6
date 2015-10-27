package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	listenAddr   = flag.String("listen", ":8080", "HTTP listen adddress.")
	pollInterval = flag.Duration("poll", time.Minute, "Polls followers every `t`.")
)

func main() {
	flag.Parse()
	s := NewServer()
	http.Handle("/", s)
	fmt.Printf("Listening for requests on %s...\n", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

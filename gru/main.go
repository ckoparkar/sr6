package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	listenAddr        = flag.String("listen", ":8281", "HTTP listen adddress.")
	pollInterval      = flag.Duration("poll", time.Minute, "Polls followers every `t`.")
	sshPrivateKeyPath = flag.String("ssh-private-key", "~/.ssh/id_rsa", "Path to SSH private key.")
	sshPublicKeyPath  = flag.String("ssh-public-key", "~/.ssh/id_rsa.pub", "Path to SSH public key.")
	// must contain ID
	hostnamePattern = flag.String("hostname", "minionID", "ID in `S` is replaced with actual ids.")
	followerPort    = flag.String("follower-port", ":8282", "Port on which all follower servers run.")
)

func main() {
	flag.Parse()
	s, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", s)
	fmt.Printf("Listening for requests on %s...\n", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

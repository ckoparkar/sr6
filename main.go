package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

var (
	nodeName = flag.String("nodename", "hostname", "Advertise this nodename.")
	mode     = flag.String("mode", "server", "Decides whether to run as client/server.")
	listen   = flag.String("listen", ":8080", "HTTP listen address.")
)

func main() {
	flag.Parse()
	c := DefaultConfig()
	s, err := NewServer(c)
	if err != nil {
		log.Fatal(err)
	}
	for {
		_, err := s.serfLAN.Join([]string{serfPort}, true)
		if err != nil {
			log.Println(err)
		} else {
			break
		}
		time.Sleep(2 * time.Second)
	}
	for {
		fmt.Println(s.serfLAN.Members())
		time.Sleep(5 * time.Second)
	}
}

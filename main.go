package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	listenAddr = flag.String("listen", ":8080", "HTTP listen adddress.")
)

func main() {
	flag.Parse()
	s := new(Server)
	http.Handle("/", s)
	fmt.Printf("Listening for requests on %s...\n", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

type Server struct {
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/register" {
		body, _ := ioutil.ReadAll(r.Body)
		var req map[string]interface{}
		json.Unmarshal(body, &req)
		fmt.Println(req)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	s := NewServer()
	http.Handle("/", s)
	fmt.Printf("Listening for requests on %d...\n", 8080)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

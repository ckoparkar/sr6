package sr6

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HttpBasicAuth is used to authenticate http client with HTTP Basic Authentication
type HttpBasicAuth struct {
	// Username to use for HTTP Basic Authentication
	Username string

	// Password to use for HTTP Basic Authentication
	Password string
}

// Request is used to help build up a request
type Request struct {
	client *http.Client
	method string
	url    *url.URL
	params url.Values
	auth   *HttpBasicAuth
	body   io.Reader
}

func NewRequest(method, scheme, address, path, auth string, body io.Reader, client *http.Client) *Request {
	r := &Request{
		method: method,
		url: &url.URL{
			Scheme: scheme,
			Host:   address,
			Path:   path,
		},
		params: make(map[string][]string),
		body:   body,
	}
	if client == nil {
		r.client = http.DefaultClient
	} else {
		r.client = client
	}

	if auth != "" {
		split := strings.SplitN(auth, ":", 2)
		r.auth = &HttpBasicAuth{Username: split[0], Password: split[1]}
	}

	return r
}

// Do runs a request with our client
func (r *Request) Do() (time.Duration, *http.Response, error) {
	req, err := r.toHTTP()
	if err != nil {
		return 0, nil, err
	}
	quit := make(chan int, 0)
	var diff time.Duration
	go func() {
		ticker := time.NewTicker(time.Second)
		log.Printf("[INFO] Running HTTP request %s", req.URL)
		for {
			select {
			case <-quit:
				fmt.Println(" Done in: ", diff)
				return
			case <-ticker.C:
				fmt.Print(".")
			}
		}
	}()
	start := time.Now()
	resp, err := r.client.Do(req)
	diff = time.Now().Sub(start)
	quit <- 1
	return diff, resp, err
}

// toHTTP converts the request to an HTTP request
func (r *Request) toHTTP() (*http.Request, error) {
	// Encode the query parameters
	r.url.RawQuery = r.params.Encode()

	// Create the HTTP request
	req, err := http.NewRequest(r.method, r.url.RequestURI(), r.body)
	if err != nil {
		return nil, err
	}
	if r.auth != nil {
		req.SetBasicAuth(r.auth.Username, r.auth.Password)
	}

	req.URL.Host = r.url.Host
	req.URL.Scheme = r.url.Scheme
	req.Host = r.url.Host

	return req, nil
}

package request

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

// Request is used to help build up a request
type Request struct {
	client *http.Client
	method string
	url    *url.URL
	params url.Values
	body   io.Reader
}

func NewRequest(method, scheme, address, path string, body io.Reader, client *http.Client) *Request {
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
	return r
}

// Do runs a request with our client
func (r *Request) Do() (time.Duration, *http.Response, error) {
	req, err := r.toHTTP()
	if err != nil {
		return 0, nil, err
	}
	start := time.Now()
	resp, err := r.client.Do(req)
	diff := time.Now().Sub(start)
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

	req.URL.Host = r.url.Host
	req.URL.Scheme = r.url.Scheme
	req.Host = r.url.Host

	return req, nil
}

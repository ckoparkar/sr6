package sr6

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var ErrAmbariCredsError = errors.New("AMBARI_HTTP_AUTH is not set.")

type AmbariConfig struct {
	// Addr is the server url:port
	Addr string

	// Auth is the ambari user:password
	Auth string
}

func DefaultAmbariConfig() (*AmbariConfig, error) {
	auth := os.Getenv("AMBARI_HTTP_AUTH")
	if auth == "" {
		return nil, ErrAmbariCredsError
	}
	return &AmbariConfig{
		Addr: "localhost:8400",
		Auth: auth,
	}, nil
}

type ambari struct {
	config *AmbariConfig
}

func (a *ambari) listClusters() error {
	req := NewRequest("GET", "http", a.config.Addr, "/api/v1/clusters", a.config.Auth, nil, nil)
	rtt, resp, err := req.Do()
	if err != nil {
		return err
	}
	fmt.Println("done in", rtt)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}

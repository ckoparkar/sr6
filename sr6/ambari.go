package sr6

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var ErrAmbariCredsError = errors.New("AMBARI_HTTP_AUTH is not set.")

type ListClustersResponse struct {
	Items []struct {
		Href     string `json:"href"`
		Clusters struct {
			Name    string `json:"cluster_name"`
			Version string `json:"version"`
		} `json:"Clusters"`
	} `json:"items"`
}

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

type Ambari struct {
	config *AmbariConfig
}

func (a *Ambari) ListClusters() (*ListClustersResponse, error) {
	req := NewRequest("GET", "http", a.config.Addr, "/api/v1/clusters", a.config.Auth, nil, nil)
	_, resp, err := req.Do()
	if err != nil {
		return nil, err
	}
	var lcr ListClustersResponse
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &lcr)
	return &lcr, nil
}

func (a *Ambari) AddHost(hostname string) error {
	clusters, err := a.ListClusters()
	if err != nil {
		return err
	}
	c := clusters.Items[0]
	addr := fmt.Sprintf("%s/hosts/%s", c.Href, hostname)
	req := NewRequest("POST", "http", addr, "", a.config.Auth, nil, nil)
	_, resp, err := req.Do()
	if err != nil || resp.StatusCode != http.StatusCreated {
		return err
	}
	return nil
}

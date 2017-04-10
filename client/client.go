/*
Package client is a Go client for the Liima API.
*/
package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

//TODO: validate config, parse url

// Client is the API client that performs all operations
// against a liima server.
type Client struct {
	url string
	// config options of the liima client
	config *Config

	// client used to send and receive http requests.
	client *http.Client
}

// NewClient creates a new liima client from the config
func NewClient(config *Config) (*Client, error) {
	tr, err := newTransport(config)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   300 * time.Second,
	}

	return &Client{
		client: httpClient,
		config: config,
		url:    config.Host,
	}, nil
}

func (c *Client) setBasicAuth(request *http.Request) {
	if c.config.Username != "" {
		request.SetBasicAuth(c.config.Username, c.config.Password)
	}
}

func newTransport(config *Config) (*http.Transport, error) {
	tlsConfg, err := newTLSClientConfig(config)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          10,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConfg,
	}
	return tr, nil
}

func newTLSClientConfig(config *Config) (*tls.Config, error) {
	var certs []tls.Certificate
	caCertPool := x509.NewCertPool()

	if config.CertFile != "" {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}

	if config.CAFile != "" {
		caCert, err := ioutil.ReadFile(config.CAFile)
		if err != nil {
			return nil, err
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}

	tlsConfig := &tls.Config{
		Certificates:       certs,
		RootCAs:            caCertPool,
		InsecureSkipVerify: config.InsecureSkipVerify,
	}
	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}

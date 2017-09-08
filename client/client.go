/*
Package client is a Go client for the Liima API.
*/
package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

//DoRequest set up a json for the given url and calls the llima client.
//Method: http.MethodX
//URL: Resturl
//The bodyType will be marshaled to the rest body, depending the method
//The result will be unmarshaled to the responseType
func (c *Client) DoRequest(method string, url string, bodyType interface{}, responseType interface{}) error {

	//Setup body if MethodPost
	bData := []byte{}
	if method == http.MethodPost {
		bDataloc, err := json.Marshal(bodyType)
		if err != nil {
			log.Fatal(err)
		}
		bData = bDataloc
	}

	var bodydata = bytes.NewBuffer(bData)

	//Setup request with format "application/json"
	//ToDo: validate config.host (ending slash)
	reqURL := fmt.Sprintf(c.config.Host + url)
	req, err := http.NewRequest(method, reqURL, bodydata)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Do request
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println("Error http request: ", reqURL)
		if resp != nil {
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
		}
		return err
	}
	defer resp.Body.Close()

	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Response Body Error on request: ", reqURL)
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		return err
	}

	//Unmarshal json respond to responseType
	return json.Unmarshal(data, responseType)

}

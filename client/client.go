/*
Package client is a Go client for the Liima API.
*/
package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
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
	errs := config.Validate()
	if len(errs) != 0 {
		msg := "Config is invalide: "
		for _, err := range errs {
			msg += fmt.Sprintf("%s; ", err)
		}
		return nil, errors.New(msg)
	}

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
		MinVersion:         tls.VersionTLS12,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
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
			return fmt.Errorf("Couldn't marshal body %v, %v", bodyType, err)
		}
		bData = bDataloc
	}

	var bodydata = bytes.NewBuffer(bData)

	//Setup request with format "application/json"
	//ToDo: validate config.host (ending slash)
	reqURL := c.config.Host + url
	req, err := http.NewRequest(method, reqURL, bodydata)
	if err != nil {
		return fmt.Errorf("Cloudn't create request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Do request, retry if fails
	var resp *http.Response
	nbrOfRetry := 3
	for i := 0; i < nbrOfRetry; i++ {

		respond, err := c.client.Do(req)
		if err != nil {
			log.Print("Error http request: ", reqURL)
			log.Print("ERROR: ", err)
			if respond != nil {
				log.Print("response Status:", respond.Status)
				log.Print("response Headers:", respond.Header)
			}
			//return if max retry reached
			if i >= (nbrOfRetry - 1) {
				return err
			}
		} else {
			resp = respond
			defer respond.Body.Close()
			break
		}
	}

	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	//Check on error
	if err != nil || !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {

		//Error response if node active=false in liima appserver configuration
		if resp.StatusCode != http.StatusFailedDependency {
			log.Print("Response Error on request: ", reqURL)
			log.Print("response Status:", resp.Status)
			log.Print("response Headers:", resp.Header)
			log.Print("response Body:", string(data))
		}
		if err != nil {
			return err
		}
		return fmt.Errorf(resp.Status + " : " + string(data))
	}

	//Unmarshal json respond to responseType
	err = json.Unmarshal(data, responseType)
	if err != nil {
		return fmt.Errorf("Couldn't unmarshal response: %s\n %s", err, data)
	}
	return nil
}

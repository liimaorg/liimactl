package client

import "fmt"

// Config options for the liima client
type Config struct {
	// URL to the base of the liima server
	Host string

	// Server requires Basic authentication
	Username string
	Password string

	// TLSClientConfig contains settings to enable transport layer security
	TLSClientConfig
}

// TLSClientConfig contains settings to enable transport layer security
type TLSClientConfig struct {
	// Server requires TLS client certificate authentication
	CertFile string
	// Server requires TLS client certificate authentication
	KeyFile string
	// Trusted root certificates for server
	CAFile string
	// InsecureSkipVerify controls whether a client verifies the
	// server's certificate chain and host name.
	// TODO: implement
	InsecureSkipVerify bool
}

// Validate the configuration
func (config *Config) Validate() []error {
	//TODO: validate Host

	return config.TLSClientConfig.Validate()
}

// Validate TLSClientConfig
func (config *TLSClientConfig) Validate() []error {
	validationErrors := make([]error, 0)

	if config.CertFile != "" && config.KeyFile == "" {
		validationErrors = append(validationErrors, fmt.Errorf("KeyFile can't be empty if CertFile is set: %v", config))
	}

	return validationErrors
}

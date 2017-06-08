# Install

```
# Install dep tool
go get -u github.com/golang/dep/cmd/dep

# Install dependencies
dep ensure

go install
```

# Configuration

limactl can be configured via commandline parameters or via config file in $HOME/.liimactl/config.yaml. Sample config:

```
Host: https://liima-host/AMW_rest/
Username: user for basic auth (optional)
Password: password for basic auth (optional)
# For client cert auth
TLSClientConfig:
    CertFile: path to public key in pem format (optional)
    KeyFile: path to unencrypted private key in pem format (optional)
    CAFile: path to ca certs in pem format (optional)
    InsecureSkipVerify: false (default false)
```

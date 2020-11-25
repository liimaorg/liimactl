![Build Status](https://github.com/liimaorg/liimactl/workflows/test/badge.svg)

# Install

```
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

# Releasing

Create a new GIT tag and push it:
```
git tag -a v0.0.2 -m "Liimactl release v0.0.2"
git push origin v0.0.2
```
Travis-ci will then build the tag, create a new release on the github page and upload the binaries (win, mac, linux). After that you can add release notes and publish the release.

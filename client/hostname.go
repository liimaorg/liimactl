package client

import (
	"fmt"
	"net/http"

	"github.com/liimaorg/liimactl/client/util"
)

//Hostnames type
type Hostnames []struct {
	Host             string `json:"host"`
	AppServer        string `json:"appServer"`
	AppServerRelease string `json:"appServerRelease"`
	Runtime          string `json:"runtime"`
	Node             string `json:"node"`
	NodeRelease      string `json:"nodeRelease"`
	Environment      string `json:"environment"`
	Domain           string `json:"domain"`
	DefinedOnNode    bool   `json:"definedOnNode"`
}

//sort.Interface
func (slice Hostnames) Len() int {
	return len(slice)
}

func (slice Hostnames) Less(i, j int) bool {
	return slice[i].Domain < slice[j].Domain
}

func (slice Hostnames) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

//CommandOptionsHostName used for the command options (flags)
type CommandOptionsHostName struct {
	AppServer    []string `json:"appServer"`
	Runtime      []string `json:"runtime"`
	Environment  []string `json:"environment"`
	Host         []string `json:"host"`
	Node         []string `json:"node"`
	DisableMerge bool     `json:"disableMerge"`
}

//GetHostname return the hostnames from the client
func GetHostname(cli *Cli, commandOptions *CommandOptionsHostName) (Hostnames, error) {

	//Build URL
	url := fmt.Sprintf("resources/./hostNames?")
	url += util.BuildCommandURL(commandOptions)

	//Call rest client
	hostnames := Hostnames{}
	if err := cli.Client.DoRequest(http.MethodGet, url, nil, &hostnames); err != nil {
		return hostnames, fmt.Errorf("Error in rest call: %v", err)
	}

	return hostnames, nil
}

package client

import (
	"net/http"
	"net/url"
)

//CommandOptionsGetFilteredDeployments used for the command options filter (flags)
type CommandOptionsGetFilteredDeployments struct {
	Filter string
}

//GetFilteredDeployments return the deployment from the client
func GetFilteredDeployments(cli *Cli, commandOptions *CommandOptionsGetFilteredDeployments) (Deployments, error) {

	//Build URL
	resturl := "resources/deployments/filter?filters="
	resturl += url.QueryEscape(commandOptions.Filter)

	//Call rest client
	deployments := Deployments{}
	if err := cli.Client.DoRequest(http.MethodGet, resturl, nil, &deployments); err != nil {
		return deployments, err
	}

	return deployments, nil
}

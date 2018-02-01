package client

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

//CommandOptionsGetDeploymentFilter used for the command options filter (flags)
type CommandOptionsGetDeploymentFilter struct {
	Filter string
}

//GetDeploymentFilter return the deployment from the client
func GetDeploymentFilter(cli *Cli, commandOptions *CommandOptionsGetDeploymentFilter) Deployments {

	//Build URL
	resturl := fmt.Sprintf("resources/deployments/filter?filters=")
	resturl += url.QueryEscape(commandOptions.Filter)

	// fmt.Println(resturl)
	// os.Exit(0)

	//Call rest client
	deployments := Deployments{}
	if err := cli.Client.DoRequest(http.MethodGet, resturl, nil, &deployments); err != nil {
		log.Fatal("Error rest call: ", err)
	}

	return deployments
}

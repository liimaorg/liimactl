package client

import (
	"fmt"
	"log"
	"net/http"

	"github.com/liimaorg/liimactl/client/util"
)

//Deployments is a collection of DeplyomentResponse
type Deployments []DeplyomentResponse

//DeplyomentResponse type is autogenerated with https://mholt.github.io/json-to-go/
type DeplyomentResponse struct {
	ID              int    `json:"id"`
	TrackingID      int    `json:"trackingId"`
	State           string `json:"state"`
	DeploymentDate  int64  `json:"deploymentDate"`
	AppServerName   string `json:"appServerName"`
	AppServerID     int    `json:"appServerId"`
	AppsWithVersion []struct {
		ApplicationName string `json:"applicationName"`
		Version         string `json:"version"`
	} `json:"appsWithVersion"`
	DeploymentParameters []interface{} `json:"deploymentParameters"`
	EnvironmentName      string        `json:"environmentName"`
	ReleaseName          string        `json:"releaseName"`
	RuntimeName          string        `json:"runtimeName"`
	RequestUser          string        `json:"requestUser"`
	ConfirmUser          string        `json:"confirmUser"`
	CancelUser           interface{}   `json:"cancelUser"`
	NodeJobs             []interface{} `json:"nodeJobs"`
	CancleUser           interface{}   `json:"cancleUser"`
}

//sort.Interface
func (slice Deployments) Len() int {
	return len(slice)
}

func (slice Deployments) Less(i, j int) bool {
	return slice[i].AppServerName < slice[j].AppServerName
}

func (slice Deployments) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

//CommandOptionsGetDeployment used for the command options (flags)
type CommandOptionsGetDeployment struct {
	AppName         []string `json:"appName"`
	AppServer       []string `json:"appServerName"`
	DeploymentState []string `json:"deploymentState"`
	Environment     []string `json:"environmentName"`
	TrackingID      int      `json:"trackingId"`
	OnlyLatest      bool     `json:"onlyLatest"`
}

//Enumeration of deplyoment state
const (
	DeploymentStateSuccess        = "success"
	DeploymentStateFailed         = "failed"
	DeploymentStateCanceled       = "canceled"
	DeploymentStateRejected       = "rejected"
	DeploymentStateReadyForDeploy = "ready_for_deploy"
	DeploymentStatePreDeploy      = "pre_deploy"
	DeploymentStateProgress       = "progress"
	DeploymentStateSimulating     = "simulating"
	DeploymentStateDelayed        = "delayed"
	DeploymentStateScheduled      = "scheduled"
	DeploymentStateRequested      = "requested"
)

//GetDeployment return the deployment from the client
func GetDeployment(cli *Cli, commandOptions *CommandOptionsGetDeployment) Deployments {

	//Build URL
	url := fmt.Sprintf("resources/./deployments?")
	url += util.BuildCommandURL(commandOptions)

	//Call rest client
	deployments := Deployments{}
	if err := cli.Client.DoRequest(http.MethodGet, url, nil, &deployments); err != nil {
		log.Fatal("Error rest call: ", err)
	}

	return deployments
}
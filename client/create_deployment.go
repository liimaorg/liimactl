package client

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"time"

	"github.com/liimaorg/liimactl/client/util"
)

//DateTimeFormat defines the actual date time format
const DateTimeFormat = "2006-01-02 15:04MST"

//LiimaDateTimeFormat defines the format for Liima UTC
const LiimaDateTimeFormat = "2006-01-02T15:04:05-0700"

//appsWithVersion type
type appsWithVersion struct {
	ApplicationName string `json:"applicationName"`
	Version         string `json:"version"`
}

//deploymentParameters type
type deploymentParameters struct {
	Value string `json:"value"`
	Key   string `json:"key"`
}

//DeplyomentRequest type
type DeplyomentRequest struct {
	ReleaseName          *string                `json:"releaseName"`
	AppServerName        string                 `json:"appServerName"`
	EnvironmentName      string                 `json:"environmentName"`
	AppsWithVersion      []appsWithVersion      `json:"appsWithVersion"`
	DeploymentParameters []deploymentParameters `json:"deploymentParameters"`
	StateToDeploy        string                 `json:"stateToDeploy"`
	ContextIds           []string               `json:"contextIds"`
	DeploymentDate       string                 `json:"deploymentDate"`
	SendEmail            bool                   `json:"sendEmail"`
	RequestOnly          bool                   `json:"requestOnly"`
	Simulate             bool                   `json:"simulate"`
	ExecuteShakedownTest bool                   `json:"executeShakedownTest"`
	NeighbourhoodTest    bool                   `json:"neighbourhoodTest"`
}

//CommandOptionsCreateDeployment used for the command options (flags)
type CommandOptionsCreateDeployment struct {
	AppServer            string   `json:"appServerName"`
	AppName              []string `json:"applicationName"`
	AppVersion           []string `json:"version"`
	Environment          string   `json:"environmentName"`
	Release              string   `json:"releaseName"`
	DeploymentDate       string   `json:"deploymentDate"`
	ExecuteShakedownTest bool     `json:"executeShakedownTest"`
	Key                  []string `json:"key"`
	Value                []string `json:"value"`
	Wait                 bool     //Wait as long the WaitTime until the deplyoment success or failed
	MaxWaitTime          int      //Max wait time [seconds] until the deplyoment success or failed
	FromEnvironment      string   //Deploy last deplyoment from given environment
}

//Validate the given command options
func (commandOption *CommandOptionsCreateDeployment) validate() error {

	//Errorlist
	var errorList []string
	//Checks and add to errorList if an error
	util.Check(&errorList, commandOption.AppServer != "", "want appServer")
	util.Check(&errorList, len(commandOption.Key) == len(commandOption.Value), "want same count of key and value, got key %d != value %d", len(commandOption.Key), len(commandOption.Value))
	util.Check(&errorList, util.ValidateSingleChar(commandOption.Environment), "want environment with one char, got %s", commandOption.Environment)
	//Copy from environment, don't check AppName and AppVersion
	if commandOption.FromEnvironment != "" {
		util.Check(&errorList, util.ValidateSingleChar(commandOption.FromEnvironment), "want FromEnvironment with one char, got %s", commandOption.FromEnvironment)
	} else {
		util.Check(&errorList, len(commandOption.AppName) > 0, "want appName")
		util.Check(&errorList, len(commandOption.AppVersion) > 0, "want appVersion")
		util.Check(&errorList, len(commandOption.AppName) == len(commandOption.AppVersion), "want same count of appName and appVersion, got appName %d != appVersion %d", len(commandOption.AppName), len(commandOption.AppVersion))
	}
	//Return all errors as one
	if len(errorList) > 0 {
		return errors.New(strings.Join(errorList, ", "))
	}
	return nil
}

//CreateDeployment create a deployment and returns the deploymentresponse from the client
func CreateDeployment(cli *Cli, commandOptions *CommandOptionsCreateDeployment) (*DeploymentResponse, error) {

	if err := commandOptions.validate(); err != nil {
		return nil, err
	}

	//Build URL
	url := fmt.Sprintf("resources/./deployments")

	//Create request (body)
	deploymentRequest := DeplyomentRequest{}
	deploymentRequest.AppServerName = commandOptions.AppServer
	deploymentRequest.EnvironmentName = commandOptions.Environment
	deploymentRequest.ExecuteShakedownTest = commandOptions.ExecuteShakedownTest
	deploymentRequest.ReleaseName = &commandOptions.Release
	if commandOptions.Release == "" {
		deploymentRequest.ReleaseName = nil
	}
	//Set deploymentdate
	actTimeZone, _ := time.Now().In(time.Local).Zone() //Load act timezone
	//Parse time in actual timezone
	t, _ := time.Parse(DateTimeFormat, commandOptions.DeploymentDate+actTimeZone)
	//Format to liima UTC format
	deploymentRequest.DeploymentDate = t.Format(LiimaDateTimeFormat)

	//Get application and version from last deployment of given "from environment"
	if commandOptions.FromEnvironment != "" {
		commandOptionsGet := CommandOptionsGetDeployment{}
		commandOptionsGet.Environment = []string{commandOptions.FromEnvironment}
		commandOptionsGet.AppServer = []string{commandOptions.AppServer}
		commandOptionsGet.TrackingID = -1
		commandOptionsGet.OnlyLatest = true
		//Get last deployment
		deployments, err := GetDeployment(cli, &commandOptionsGet)
		if err != nil {
			return nil, err
		}
		if len(deployments) == 0 {
			return nil, fmt.Errorf("There was an error on creating the deplyoment, no deployment found from environment: %s", commandOptions.FromEnvironment)
		}
		lastDeployment := deployments[0]
		//Set app and version
		for i := 0; i < len(lastDeployment.AppsWithVersion); i++ {
			appVersion := appsWithVersion{
				ApplicationName: lastDeployment.AppsWithVersion[i].ApplicationName,
				Version:         lastDeployment.AppsWithVersion[i].Version,
			}
			deploymentRequest.AppsWithVersion = append(deploymentRequest.AppsWithVersion, appVersion)
		}
	} else {
		//Application and version
		for i := 0; i < len(commandOptions.AppName); i++ {
			appVersion := appsWithVersion{
				ApplicationName: commandOptions.AppName[i],
				Version:         commandOptions.AppVersion[i],
			}
			deploymentRequest.AppsWithVersion = append(deploymentRequest.AppsWithVersion, appVersion)
		}
	}
	//Deployment parameter
	for i := 0; i < len(commandOptions.Key); i++ {
		dpKey := deploymentParameters{
			Key:   commandOptions.Key[i],
			Value: commandOptions.Value[i],
		}
		deploymentRequest.DeploymentParameters = append(deploymentRequest.DeploymentParameters, dpKey)
	}

	//Call rest client
	deploymentResponse := &DeploymentResponse{}
	if err := cli.Client.DoRequest(http.MethodPost, url, &deploymentRequest, &deploymentResponse); err != nil {
		return deploymentResponse, err
	}

	//Wait on deplyoment success or failed
	if commandOptions.Wait && commandOptions.MaxWaitTime > 5 {
		commandOptionsGet := CommandOptionsGetDeployment{
			TrackingID: deploymentResponse.TrackingID,
		}

		//Timeout 10min = 600sec / 5sec = 120 counts
		maxCounts := commandOptions.MaxWaitTime / 5
		for i := 0; i < maxCounts; i++ {
			deployments, err := GetDeployment(cli, &commandOptionsGet)
			if err != nil {
				return nil, err
			}

			if len(deployments) != 1 {
				return nil, fmt.Errorf("There was an error on creating the deplyoment, no deployment get")
			}

			log.Printf("AppServer: %-30s State: %-20s\n", deployments[0].AppServerName, deployments[0].State)

			deploymentResponse = &deployments[0]
			if deployments[0].State == DeploymentStateFailed || deployments[0].State == DeploymentStateSuccess {
				break
			}
			if i < maxCounts-1 {
				time.Sleep(time.Second * 5)
			} else {
				return nil, fmt.Errorf("Timeout on deployment")
			}
		}
	}

	//Return response
	return deploymentResponse, nil
}

package client

import (
	"errors"
	"log"
	"strings"

	"github.com/liimaorg/liimactl/client/util"
)

//CommandOptionsPromoteDeployments used for the command options (flags)
type CommandOptionsPromoteDeployments struct {
	Environment          string   `json:"environmentName"`
	DeploymentDate       string   `json:"deploymentDate"`
	ExecuteShakedownTest bool     `json:"executeShakedownTest"`
	Wait                 bool     //Wait as long the WaitTime until the deplyoment success or failed
	MaxWaitTime          int      //Max wait time [seconds] until the deplyoment success or failed
	FromEnvironment      string   //Deploy last deplyoment from given environment
	BlacklistAppServer   []string //Blacklist with all appServer, which should not be deployed
	BlacklistRuntime     []string //Blacklist with all runtimes, which should not be deployed
	Silent               bool     //silent mode, no confirmation of promote the whole environment
}

//Validate the given command options
func (commandOption *CommandOptionsPromoteDeployments) validate() error {

	//Errorlist
	var errorList []string
	//Checks and add to errorList if an error
	util.Check(&errorList, util.ValidateSingleChar(commandOption.Environment), "want environment with one char, got %s", commandOption.Environment)
	util.Check(&errorList, util.ValidateSingleChar(commandOption.FromEnvironment), "want FromEnvironment with one char, got %s", commandOption.FromEnvironment)

	//Return all errors as one
	if len(errorList) > 0 {
		return errors.New(strings.Join(errorList, ", "))
	}
	return nil
}

//PromoteDeployments creates multiple deployments and returns the deploymentresponse
func PromoteDeployments(cli *Cli, commandOptions *CommandOptionsPromoteDeployments) (Deployments, error) {

	//validate commandoptions
	if err := commandOptions.validate(); err != nil {
		log.Println("Error command validation: ", err)
		return Deployments{}, err
	}

	//Create the filter for searching all deplyoments from an environment
	commandOptionsGetFilter := CommandOptionsGetFilteredDeployments{}
	filter := []string{`[{"name":"Environment","comp":"eq","val":"`, commandOptions.FromEnvironment, `"},{"name":"Latest deployment job for App Server and Env","comp":"eq","val":"true"}]`}
	commandOptionsGetFilter.Filter = strings.Join(filter, "")

	//Get all last deployments of the given environment
	deployments, err := GetFilteredDeployments(cli, &commandOptionsGetFilter)
	if err != nil {
		log.Println("Error on getting the filtered deployments: ", err)
		return deployments, err
	}
	if len(deployments) == 0 {
		log.Println("There was an error on creating the deplyoment, no deployment found from environment: ", commandOptions.FromEnvironment)
		return deployments, err
	}

	//Remove all deployments
	// - defined in the blacklist
	// - with a runtime on the blacklistRuntime
	// - without state "success"
	for i := len(deployments) - 1; i >= 0; i-- {
		if (util.Contains(deployments[i].AppServerName, commandOptions.BlacklistAppServer)) || (deployments[i].State != DeploymentStateSuccess) || util.Contains(deployments[i].RuntimeName, commandOptions.BlacklistRuntime) {
			deployments = append(deployments[:i], deployments[i+1:]...)
		}
	}

	//Create deployments
	createdDeployments := Deployments{}

	for _, actDeployment := range deployments {

		commandOptionsCreateDeployment := CommandOptionsCreateDeployment{}
		commandOptionsCreateDeployment.AppServer = actDeployment.AppServerName
		commandOptionsCreateDeployment.Release = actDeployment.ReleaseName
		commandOptionsCreateDeployment.Environment = commandOptions.Environment
		commandOptionsCreateDeployment.DeploymentDate = commandOptions.DeploymentDate
		commandOptionsCreateDeployment.AppName = make([]string, len(actDeployment.AppsWithVersion))
		commandOptionsCreateDeployment.AppVersion = make([]string, len(actDeployment.AppsWithVersion))
		for i := range actDeployment.AppsWithVersion {
			commandOptionsCreateDeployment.AppName[i] = actDeployment.AppsWithVersion[i].ApplicationName
			commandOptionsCreateDeployment.AppVersion[i] = actDeployment.AppsWithVersion[i].Version
		}

		deplyoment, err := CreateDeployment(cli, &commandOptionsCreateDeployment)
		if err != nil {
			log.Println("Error Create Deployment: ", err)
			return createdDeployments, err
		}
		createdDeployments = append(createdDeployments, *deplyoment)
	}

	//Return response
	return createdDeployments, nil
}

package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/liimaorg/liimactl/client/util"
)

//CommandOptionsInstallDeployment used for the command options (flags)
type CommandOptionsInstallDeployment struct {
	Environment          string `json:"environmentName"`
	DeploymentDate       string `json:"deploymentDate"`
	ExecuteShakedownTest bool   `json:"executeShakedownTest"`
	Wait                 bool   //Wait as long the WaitTime until the deplyoment success or failed
	MaxWaitTime          int    //Max wait time [seconds] until the deplyoment success or failed
	FromEnvironment      string //Deploy last deplyoment from given environment
	Blacklist            string //Blacklist with all appServer, which should not be deployed [filepath]
}

//Validate the given command options
func (commandOption *CommandOptionsInstallDeployment) validate() error {

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

//InstallDeployment create multiple deployments and returns the deploymentresponse
func InstallDeployment(cli *Cli, commandOptions *CommandOptionsInstallDeployment) (Deployments, error) {

	//validate commandoptions
	if err := commandOptions.validate(); err != nil {
		log.Fatal("Error command validation: ", err)
		return Deployments{}, err
	}

	//Create the filter for searching all deplyoments from an environment
	commandOptionsGetFilter := CommandOptionsGetDeploymentFilter{}
	filter := []string{`[{"name":"Environment","comp":"eq","val":"`, commandOptions.FromEnvironment, `"},{"name":"Latest deployment job for App Server and Env","comp":"eq","val":"true"}]`}
	commandOptionsGetFilter.Filter = strings.Join(filter, "")

	//Get all last deployments of the given environment
	deployments := GetDeploymentFilter(cli, &commandOptionsGetFilter)
	if len(deployments) == 0 {
		log.Fatal("There was an error on creating the deplyoment, no deployment found from environment: ", commandOptions.FromEnvironment)
	}

	//Build array server (blacklist)
	blacklistServers := []string{}
	if commandOptions.Blacklist != "" {
		content, err := ioutil.ReadFile(commandOptions.Blacklist)
		if err != nil {
			log.Fatal("Error reading file blacklist: ", err)
		}
		blacklistServers = strings.Split(string(content), "\n")
	}

	//Array of runtime (blacklist)
	blacklistRuntime := []string{"Kubernetes", "Kube_helm"}

	//Remove all deployments defined in the blacklist and remove all deployments without state "success" and all deployments w
	for i := len(deployments) - 1; i >= 0; i-- {
		if (util.Contains(deployments[i].AppServerName, blacklistServers)) || (deployments[i].State != DeploymentStateSuccess) || util.Contains(deployments[i].RuntimeName, blacklistRuntime) {
			deployments = append(deployments[:i], deployments[i+1:]...)
		}
	}

	//Create deployments
	createdDeployments := Deployments{}

	//Ask user for confirmation
	msg := fmt.Sprintf("Do you really want to start the deployment of %d app-servers on environment: %s", len(deployments), commandOptions.Environment)
	if util.AskYesNo(msg) {

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
				log.Fatal("Error Create Deployment: ", err)
			}
			createdDeployments = append(createdDeployments, deplyoment)
		}
	}

	//Return response
	return createdDeployments, nil
}

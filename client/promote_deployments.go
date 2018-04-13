package client

import (
	"errors"
	"log"
	"strings"
	"time"

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
	WhitelistAppServer   []string //Whitelist with all appServer, which should be deployed, if not WhitelistAppServer is defined, the whole environment will deployed (exclusive blacklist)
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
		return nil, err
	}

	//Create the filter for searching all deplyoments from an environment
	commandOptionsGetFilter := CommandOptionsGetDeployment{}
	commandOptionsGetFilter.Environment = []string{commandOptions.FromEnvironment}
	commandOptionsGetFilter.OnlyLatest = true
	commandOptionsGetFilter.TrackingID = -1
	commandOptionsGetFilter.AppServer = commandOptions.WhitelistAppServer

	//Get all last deployments of the given environment
	deployments, err := GetDeployment(cli, &commandOptionsGetFilter)
	if err != nil {
		log.Println("Error on getting the filtered deployments: ", err)
		return nil, err
	}

	if len(deployments) == 0 {
		log.Println("There was an error on creating the deplyoment, no deployment found from environment: ", commandOptions.FromEnvironment)
		return nil, err
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

		deployment, err := CreateDeployment(cli, &commandOptionsCreateDeployment)
		if err != nil {
			log.Println("Error Create Deployment: ", err)
			return createdDeployments, err
		}
		createdDeployments = append(createdDeployments, *deployment)
	}

	//Wait on deplyoment success or failed
	sleepTime := 60 //seconds, polling each x seconds
	if commandOptions.Wait && commandOptions.MaxWaitTime > sleepTime {

		commandOptionsGet := CommandOptionsGetDeployment{}
		commandOptionsGet.Environment = []string{commandOptions.Environment}
		commandOptionsGet.OnlyLatest = true
		commandOptionsGet.TrackingID = -1
		for _, actDeployment := range createdDeployments {
			commandOptionsGet.AppServer = append(commandOptionsGet.AppServer, actDeployment.AppServerName)
		}

		//Timeout 10min = 600sec / 30sec = 20 counts
		maxCounts := commandOptions.MaxWaitTime / sleepTime
		for i := 0; i < maxCounts; i++ {
			deployments, err := GetDeployment(cli, &commandOptionsGet)
			if err != nil {
				return nil, err
			}

			if len(deployments) == 0 {
				log.Println("There was an error on creating the deplyoment, no deployment found from environment: ", commandOptions.FromEnvironment)
				return nil, err
			}

			createdDeployments = deployments

			// Check if all deployments are finished
			allfinished := true
			for _, actDeployment := range deployments {
				log.Printf("AppServer: %-30s State: %-20s\n", actDeployment.AppServerName, actDeployment.State)
				allfinished = allfinished && (actDeployment.State == DeploymentStateFailed || actDeployment.State == DeploymentStateSuccess)
			}
			// Break loop if finished
			if allfinished {
				break
			}

			if i < maxCounts-1 {
				time.Sleep(time.Second * time.Duration(sleepTime))
			} else {
				log.Println("Timeout on promote deployment on environment: ", commandOptions.FromEnvironment)
				return nil, err
			}
		}
	}

	//Return response
	return createdDeployments, nil
}

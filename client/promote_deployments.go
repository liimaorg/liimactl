package client

import (
	"errors"
	"fmt"
	"log"
	"math"
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

//checkDeploymentResults waits for deployments to finish or maxWaitTime is reached
func checkDeploymentResults(cli *Cli, commandOptionsGet *CommandOptionsGetDeployment, maxWaitTime int) (Deployments, error) {

	checkedDeployments := Deployments{}

	const sleepTime = 60 //seconds, polling each x seconds
	//Timeout 10min = 600sec / 60sec = 10 counts
	maxCounts := int(math.Max(1, float64(maxWaitTime/sleepTime)))
	for i := 0; i <= maxCounts; i++ {

		deployments, err := GetDeployment(cli, commandOptionsGet)
		if err != nil {
			return nil, err
		}

		if len(deployments) == 0 {
			return nil, fmt.Errorf("There was an error on checking the deplyoments, no deployment found")
		}

		checkedDeployments = deployments

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
		//Check iterations, sleep or timeout
		if i < maxCounts {
			time.Sleep(time.Second * time.Duration(sleepTime))
		} else {
			return checkedDeployments, fmt.Errorf("Timeout on checking deployment results")
		}
	}

	return checkedDeployments, nil
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
		return deployments, err
	}

	if len(deployments) == 0 {
		return nil, fmt.Errorf("There was an error on creating the deplyoment, no deployment found from environment: %s ", commandOptions.FromEnvironment)
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
			log.Printf("Error Create Deployment for app server: %s error: %s", actDeployment.AppServerName, err)
			return createdDeployments, err
		}
		createdDeployments = append(createdDeployments, *deployment)
	}

	//Wait on deplyoment success or failed
	if commandOptions.Wait {

		//Create filter for created deplyoments
		commandOptionsGetFilter.Environment = []string{commandOptions.Environment}
		commandOptionsGetFilter.AppServer = nil
		for _, actDeployment := range createdDeployments {
			commandOptionsGetFilter.AppServer = append(commandOptionsGetFilter.AppServer, actDeployment.AppServerName)
		}
		//Check deployments
		deployments, err := checkDeploymentResults(cli, &commandOptionsGetFilter, commandOptions.MaxWaitTime)
		if err != nil {
			return deployments, err
		}
		createdDeployments = deployments
	}

	//Return response
	return createdDeployments, nil
}

package deployment

import (
	"fmt"
	"sort"
	"time"

	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/cobra"
)

var (
	//Long command description
	deploymentGetLong = fmt.Sprintf(` 
    Get deployment with the use of specific filters.`)

	//Example command description
	deploymentGetExample = fmt.Sprintf(` 
    # Get a deplyoment with specific filters. 
    liimactl.exe deployment get --appServer=test_application --environment=I`)

	//Flags of the command
	commandOptions client.CommandOptionsDeployment
)

//newGetCommand is a command to get deployments
func newGetCommand(cli *client.Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "get [flags] ",
		Short:   "Get deployments",
		Long:    deploymentGetLong,
		Example: deploymentGetExample,
		Run: func(cmd *cobra.Command, args []string) {
			runGet(cmd, cli, args)
		},
	}

	cmd.Flags().StringSliceVarP(&commandOptions.AppName, "appName", "", []string{}, "Application Name")
	cmd.Flags().StringSliceVarP(&commandOptions.AppServer, "appServer", "", []string{}, "Application Server Name")
	cmd.Flags().StringSliceVarP(&commandOptions.DeploymentState, "deploymentState", "", []string{}, "Deplyoment State")
	cmd.Flags().StringSliceVarP(&commandOptions.Environment, "environment", "", []string{}, "	Environment Filter")
	cmd.Flags().BoolVarP(&commandOptions.OnlyLatest, "onlyLatest", "", false, "only Latest Filter")

	//CommandOptionsDeployment used for the command options (flags)
	type CommandOptionsDeployment struct {
		AppName         []string `json:"appName"`
		AppServer       []string `json:"appServerName"`
		DeploymentState []string `json:"deploymentState"`
		Environment     []string `json:"environmentName"`
		OnlyLatest      bool     `json:"onlyLatest"`
	}

	return cmd
}

//Get the deployments properties given by the arguments (see type Deplyoments) and print it on the console
func runGet(cmd *cobra.Command, cli *client.Cli, args []string) {

	//Get deplyoments
	deplyoments := client.GetDeployment(cli, &commandOptions)

	const createdFormat = "2017-07-06 21:00" //"Jan 2, 2006 at 3:04pm (MST)"

	//Print result
	sort.Sort(deplyoments)
	for _, deplyoment := range deplyoments {
		cmd.Println("------")
		cmd.Printf("%s ", deplyoment.AppServerName)
		cmd.Printf("%s ", deplyoment.EnvironmentName)
		cmd.Printf("%s ", deplyoment.ReleaseName)
		cmd.Printf("%s ", time.Unix(0, deplyoment.DeploymentDate*int64(time.Millisecond)).Format("2006-01-02T15:04"))
		cmd.Println(deplyoment.State)

		for _, appsWithVersion := range deplyoment.AppsWithVersion {
			cmd.Printf("%s ", appsWithVersion.ApplicationName)
			cmd.Println(appsWithVersion.Version)
		}
	}

}

package deployment

import (
	"fmt"
	"log"
	"sort"

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
	commandOptionsGet client.CommandOptionsGetDeployment
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

	cmd.Flags().StringSliceVarP(&commandOptionsGet.AppName, "appName", "n", []string{}, "Application Name")
	cmd.Flags().StringSliceVarP(&commandOptionsGet.AppServer, "appServer", "a", []string{}, "Application Server Name")
	cmd.Flags().StringSliceVarP(&commandOptionsGet.DeploymentState, "deploymentState", "d", []string{}, "Deplyoment State")
	cmd.Flags().StringSliceVarP(&commandOptionsGet.Environment, "environment", "e", []string{}, "Environment Filter")
	cmd.Flags().BoolVarP(&commandOptionsGet.OnlyLatest, "onlyLatest", "l", false, "Only Latest Filter")
	cmd.Flags().IntVarP(&commandOptionsGet.TrackingID, "trackingId", "t", -1, "Tracking ID")

	return cmd
}

//Get the deployments properties given by the arguments (see type Deplyoments) and print it on the console
func runGet(cmd *cobra.Command, cli *client.Cli, args []string) {

	//Get deployments
	deployments, err := client.GetDeployment(cli, &commandOptionsGet)
	if err != nil {
		log.Fatal(err)
	}

	//Print result
	sort.Sort(deployments)
	for _, deployment := range deployments {
		PrintDeployment(cmd, &deployment)
	}

}

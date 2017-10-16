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
	cmd.Flags().StringSliceVarP(&commandOptionsGet.Environment, "environment", "e", []string{}, "	Environment Filter")
	cmd.Flags().BoolVarP(&commandOptionsGet.OnlyLatest, "onlyLatest", "l", false, "only Latest Filter")
	cmd.Flags().IntVarP(&commandOptionsGet.TrackingID, "trackingId", "t", -1, "Tracking ID")

	return cmd
}

//Get the deployments properties given by the arguments (see type Deplyoments) and print it on the console
func runGet(cmd *cobra.Command, cli *client.Cli, args []string) {

	//Get deployments
	deployments := client.GetDeployment(cli, &commandOptionsGet)

	const createdFormat = "2017-07-06 21:00" //"Jan 2, 2006 at 3:04pm (MST)"

	//Print result
	sort.Sort(deployments)
	for _, deployment := range deployments {
		cmd.Println("------")
		cmd.Printf("%s ", deployment.AppServerName)
		cmd.Printf("%s ", deployment.EnvironmentName)
		cmd.Printf("%s ", deployment.ReleaseName)
		cmd.Printf("%s ", time.Unix(0, deployment.DeploymentDate*int64(time.Millisecond)).Format("2006-01-02T15:04"))
		cmd.Println(deployment.State)

		for _, appsWithVersion := range deployment.AppsWithVersion {
			cmd.Printf("%s ", appsWithVersion.ApplicationName)
			cmd.Println(appsWithVersion.Version)
		}
	}

}

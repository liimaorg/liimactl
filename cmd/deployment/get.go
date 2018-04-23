package deployment

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/cobra"
)

var (
	//Long command description
	deploymentGetLong = `	Get deployment with the use of specific filters.`

	//Example command description
	deploymentGetExample = `	# Get a deplyoment with specific filters. 
	liimactl.exe deployment get --appServer=test_application --environment=I
	# Filters can also be passed as JSON
	liimactl.exe deployment get --filter='[{"name":"Environment","comp":"eq","val":"Y"},{"name":"Application server","comp":"eq","val":"liima"}]'
	liimactl.exe deployment get --filter='[{"name":"Environment","comp":"eq","val":"Y"},{"name":"Latest deployment job for App Server and Env","comp":"eq","val":"true"}]'
	`
	//Flags of the command
	commandOptionsGet client.CommandOptionsGetDeployment
	deploymentFilter  string
	deploymentState   *[]string
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

	deploymentState = &[]string{}
	cmd.Flags().StringSliceVarP(&commandOptionsGet.AppName, "appName", "n", []string{}, "Application Name")
	cmd.Flags().StringSliceVarP(&commandOptionsGet.AppServer, "appServer", "a", []string{}, "Application Server Name")
	cmd.Flags().StringSliceVarP(deploymentState, "deploymentState", "d", []string{}, "Deplyoment State")
	cmd.Flags().StringSliceVarP(&commandOptionsGet.Environment, "environment", "e", []string{}, "Environment Filter")
	cmd.Flags().BoolVarP(&commandOptionsGet.OnlyLatest, "onlyLatest", "l", false, "Only Latest Filter")
	cmd.Flags().IntVarP(&commandOptionsGet.TrackingID, "trackingId", "t", -1, "Tracking ID")
	cmd.Flags().IntSliceVarP(&commandOptionsGet.ID, "id", "i", []int{}, "Deployment ID")
	cmd.Flags().StringVarP(&deploymentFilter, "filter", "f", "", "Deployment filter in JSON")

	return cmd
}

//Get the deployments properties given by the arguments (see type Deplyoments) and print it on the console
func runGet(cmd *cobra.Command, cli *client.Cli, args []string) {
	// convert to client types
	for _, state := range *deploymentState {
		commandOptionsGet.DeploymentState = append(commandOptionsGet.DeploymentState, client.DeploymentState(state))
	}
	if deploymentFilter != "" {
		err := json.Unmarshal([]byte(deploymentFilter), &commandOptionsGet.Filter)
		if err != nil {
			log.Fatalf("Filter is not valid: %v", err)
		}
	}

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

package deployment

import (
	"fmt"
	"sort"

	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/cobra"
)

var (
	//Long command description
	deploymentGetFilterLong = fmt.Sprintf(` 
    Get deployment with the use of given filters.`)

	//Example command description
	deploymentGetFilterExample = fmt.Sprintf(` 
    # Get a deplyoment with specific filters. 
	liimactl.exe deployment getFilter --filter='[{"name":"Environment","comp":"eq","val":"Y"},{"name":"Application server","comp":"eq","val":"aps_bau_kube"}]'
	liimactl.exe deployment getFilter --filter='[{"name":"Environment","comp":"eq","val":"Y"},{"name":"Latest deployment job for App Server and Env","comp":"eq","val":"true"}]'`)

	//Flags of the command
	commandOptionsGetFilter client.CommandOptionsGetDeploymentFilter
)

//newGetFilterCommand is a command to get deployments filter
func newGetFilterCommand(cli *client.Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "getFilter [flags] ",
		Short:   "Get deployments filter",
		Long:    deploymentGetFilterLong,
		Example: deploymentGetFilterExample,
		Run: func(cmd *cobra.Command, args []string) {
			runGetFilter(cmd, cli, args)
		},
	}

	cmd.Flags().StringVarP(&commandOptionsGetFilter.Filter, "filter", "f", "", "Filter")

	return cmd
}

//Get the deployments properties given by the arguments (see type Deplyoments) and print it on the console
func runGetFilter(cmd *cobra.Command, cli *client.Cli, args []string) {

	//Get deployments
	deployments := client.GetDeploymentFilter(cli, &commandOptionsGetFilter)

	const createdFormat = "2017-07-06 21:00" //"Jan 2, 2006 at 3:04pm (MST)"

	//Print result
	sort.Sort(deployments)
	for _, deployment := range deployments {
		PrintDeployment(cmd, deployment)
	}

}
